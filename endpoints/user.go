package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
	"golang.org/x/crypto/bcrypt"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//GetProfile returns the profile of the currently logged in user to the client
func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(user)
}

//GetProfilePicture returns the rofile icture of the UserID
func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	var userID int64
	if len(strings.Split(r.URL.String(), "/")) != 5 {
		userID, _ = strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	} else {
		user, err := utils.GetUser(r, w)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, "could not open avatar File "+err.Error())
		return
	}
	_, err = io.Copy(file, tempFile)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, "could not write avatar File "+err.Error())
	}

	w.WriteHeader(200)
}

//RemoveProfilePicture resets the profile picture
func RemoveProfilePicture(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
		dev.ReportUserError(w, "Nothing changed!")
		return
	}

	if db.PatchUser(user) != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if len(hashedPassword) > 0 {
		if db.UpdatePassword(user.ID, hashedPassword) != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	perms, err := perms.CombinePermissions(user)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
	}

	newUser.ID, err = db.CreateUser(newUser, req.Password)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(newUser)
}

//GetUsers returns all groups
func GetUsers(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	perms, err := perms.CombinePermissions(user)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(users)
}
