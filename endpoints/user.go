package endpoints

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/redis"
	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
	"golang.org/x/crypto/bcrypt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
)

//GetProfile returns the profile of the currently logged in user to the client
func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if err := db.GetSettings(&user); err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(user)
}

//GetPermissionsOfUser returns all permissions of the user
func GetPermissionsOfUser(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	req, found, err := db.GetUser(userID)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Requested user doesn't exist")
		return
	}

	userperms, err := perms.CombinePermissions(user)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !userperms.Admin && !userperms.CanChangePermissionsGlobal && req.ID != user.ID {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to view permissions in this context.")
		return
	}

	reqperms, err := perms.CombinePermissions(req)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(reqperms)
}

//PatchSettings updates the users preferences. WARNING! The full settings struct must be given to the backend in order to update the settings
func PatchSettings(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	var newSettings models.UserSettings
	rawBytes, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawBytes, &newSettings)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed: "+err.Error())
		return
	}

	if newSettings.Notification.MailNotificationFrequency < 0 {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Notification Frequency can't be negative")
		return
	}

	user.Settings = newSettings
	if err := db.PatchSettings(user); err != nil {

		utils.ReportInternalErrorToUser(err, w)
	} else {
		json.NewEncoder(w).Encode(user)
	}
}

//GetProfilePicture returns the rofile icture of the UserID
func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	var userID int64
	if len(strings.Split(r.URL.String(), "/")) != 5 {
		userID, _ = strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	} else {
		user, err := utils.GetUser(r, w)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}
		userID = user.ID
	}

	seperator := string(os.PathSeparator)
	cwd, _ := os.Getwd()
	filepath := fmt.Sprint(cwd, seperator, "userData", seperator, "avatar", seperator)
	found := false
	for _, k := range models.GetAllowedImageFormates() {
		_, err := os.Stat(filepath + strconv.FormatInt(userID, 10) + "." + k)
		if err == nil {
			found = true
			filepath += strconv.FormatInt(userID, 10) + "." + k
		}
	}
	if !found {
		filepath = fmt.Sprint(cwd, seperator, "userData", seperator, "sampleData", seperator, "defaultProfilePicture.png")
	}

	http.ServeFile(w, r, filepath)
}

//UpdateProfilePicture deletes the old and adds a new profile picture
func UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	r.ParseMultipartForm(32 << 20)
	tempFile, handler, err := r.FormFile("avatar")

	if err != nil {
		w.WriteHeader(406)
		dev.ReportUserError(w, "No 'avatar' data found")
		return
	}

	defer tempFile.Close()

	fileType := ""
	for _, k := range models.GetAllowedImageFormates() {
		if strings.HasSuffix(handler.Filename, k) {
			fileType = k
		}
	}
	if len(fileType) == 0 {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Invalid avatar filetype")
		return
	}

	seperator := string(os.PathSeparator)
	filepath, _ := os.Getwd()
	filepath += fmt.Sprint(seperator, "userData", seperator, "avatar", seperator)
	for _, k := range models.GetAllowedImageFormates() {
		_, err := os.Stat(filepath + strconv.FormatInt(user.ID, 10) + "." + k)
		if err == nil {
			os.Remove(filepath + strconv.FormatInt(user.ID, 10) + "." + k)
			break
		}
	}

	file, err := os.Create(filepath + strconv.FormatInt(user.ID, 10) + "." + fileType)
	if err != nil {

		dev.ReportError(err, w, "could not open avatar File "+err.Error())
		return
	}
	_, err = io.Copy(file, tempFile)
	if err != nil {

		dev.ReportError(err, w, "could not write avatar File "+err.Error())
	}

	w.WriteHeader(200)
}

//RemoveProfilePicture resets the profile picture
func RemoveProfilePicture(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	seperator := string(os.PathSeparator)
	filepath, _ := os.Getwd()
	filepath += fmt.Sprint(seperator, "userData", seperator, "avatar", seperator)
	for _, k := range models.GetAllowedImageFormates() {
		_, err := os.Stat(filepath + strconv.FormatInt(user.ID, 10) + "." + k)
		if err == nil {
			os.Remove(filepath + strconv.FormatInt(user.ID, 10) + "." + k)
			break
		}
	}

	w.WriteHeader(200)
}

type profileWebRequest struct {
	Username, Firstname, Lastname, Mail, Password string
}

//PatchProfile updates profile information
func PatchProfile(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	var req profileWebRequest

	rawBytes, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed")
		return
	}

	somethingChanged := false
	var hashedPassword string

	if !utils.IsEmpty(req.Username) && req.Username != user.Username {
		user.Username = req.Username
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Firstname) && req.Firstname != user.Firstname {
		user.Firstname = req.Firstname
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Lastname) && req.Lastname != user.Lastname {
		user.Lastname = req.Lastname
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Mail) && req.Mail != user.Mail {
		user.Mail = req.Mail
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Password) {
		rawPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)
		hashedPassword = string(rawPasswordBytes)
		somethingChanged = true
	}

	if !somethingChanged {
		w.WriteHeader(406)
		dev.ReportUserError(w, templates.NothingChanged)
		return
	}

	if db.PatchUser(user) != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if len(hashedPassword) > 0 {
		if db.UpdatePassword(user.ID, hashedPassword) != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}
	}

	json.NewEncoder(w).Encode(user)
}

type userWebRequest struct {
	Username, Mail, Firstname, Lastname, Password string
}

//CreateUser updates profile information
func CreateUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	perms, err := perms.CombinePermissions(user)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !perms.CanCreateUsers && !perms.Admin {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed create users")
		return
	}

	var req userWebRequest

	rawBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	err = json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Request is malformed: "+err.Error())
		return
	}

	if utils.IsEmpty(req.Password) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Password can't be empty for new user")
		return
	}

	if utils.IsEmpty(req.Username) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Username can't be empty")
		return
	}

	if utils.IsEmpty(req.Firstname) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Firstname can't be empty")
		return
	}

	if utils.IsEmpty(req.Lastname) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Lastname can't be empty")
		return
	}

	if utils.IsEmpty(req.Mail) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Mail can't be empty")
		return
	}

	users, err := db.GetALLUsers()
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	for _, k := range users {
		if strings.ToLower(k.Username) == strings.ToLower(req.Username) {
			w.WriteHeader(406)
			dev.ReportUserError(w, "User with that username already exists")
			return
		}

		if strings.ToLower(k.Mail) == strings.ToLower(req.Mail) {
			w.WriteHeader(406)
			dev.ReportUserError(w, "User with that E-Mail already exists")
			return
		}
	}

	var newUser models.User = models.User{
		Username:  req.Username,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
		Mail:      req.Mail,
		Settings: models.UserSettings{
			Notification: models.NotificationSettings{
				MailNotificationEnabled:         true,
				MailNotificationAboutNewTickets: false,
				MailNotificationAboutUpdates:    true,
				MailNotificationFrequency:       300,
			},
		},
	}

	newUser.ID, err = db.CreateUser(newUser, req.Password)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(newUser)
}

//GetUsers returns all groups
func GetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	perms, err := perms.CombinePermissions(user)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	//All Perms that allow access to all users on instance
	if !perms.Admin && !perms.CanChangePermissionsGlobal && !perms.CanModifyUsers {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to view all users")
		return
	}

	users, err := db.GetALLUsers()
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(users)
}

//GetSpecificUser returns the requested user via JSON
func GetSpecificUser(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	userid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

	user, found, err := db.GetUser(userid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Specified user wasn't found")
		return
	}

	json.NewEncoder(w).Encode(user)
}

//DeactivateUser sets the deactivated flag for the specified user
func DeactivateUser(w http.ResponseWriter, r *http.Request) {
	userid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	req, found, err := db.GetUser(userid)

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Requested User wasn't found")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	userperms, err := perms.CombinePermissions(user)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if req.ID != user.ID && !userperms.Admin {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to delete users")
		return
	}

	//Check if this is the last Admin
	if userperms.Admin {
		users, err := db.GetALLUsers()
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		var remainingAdmins int64 = 0
		for _, k := range users {
			if k.ID != req.ID {
				tempperms, err := perms.CombinePermissions(k)
				if err != nil {

					utils.ReportInternalErrorToUser(err, w)
					return
				}

				if tempperms.Admin {
					remainingAdmins++
				}
			}
		}

		if remainingAdmins == 0 {
			w.WriteHeader(406)
			dev.ReportUserError(w, "You cannot delete the last remaining admin user")
			return
		}
	}

	if err := db.DeactivateUser(req.ID); err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	RemoveProfilePicture(w, r)

	data, _ := json.Marshal(req)
	taskid, err := db.CreateTask(models.DeleteUser, string(data), sql.NullInt32{Valid: false}, sql.NullString{Valid: false}, sql.NullInt64{Valid: false}, sql.NullInt32{Valid: false})
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	redis.DestroyAllSessions(req.ID)

	json.NewEncoder(w).Encode(struct {
		TaskID int64
	}{
		taskid,
	})
}
