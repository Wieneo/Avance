package endpoints

import (
	"encoding/json"
<<<<<<< HEAD
	"fmt"
	"io"
=======
	"io/ioutil"
>>>>>>> Profiles can now be updated via PATCH /api/v1/profile
	"net/http"
	"os"
	"strconv"
	"strings"

<<<<<<< HEAD
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
=======
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
>>>>>>> Profiles can now be updated via PATCH /api/v1/profile
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
	filepath, _ := os.Getwd()
	filepath += fmt.Sprint(seperator, "userData", seperator, "avatar", seperator)
	found := false
	for _, k := range models.GetAllowedImageFormates() {
		_, err := os.Stat(filepath + strconv.FormatInt(userID, 10) + "." + k)
		if err == nil {
			found = true
			filepath += strconv.FormatInt(userID, 10) + "." + k
		}
	}
	if !found {
		filepath += "default.png"
	}

	http.ServeFile(w, r, filepath)
}

//UpdateProfilePicture sets the profile picture of the current user
func UpdateProfilePicture(w http.ResponseWriter, r *http.Request) {
type profileWebRequest struct {
	Username, Firstname, Lastname, Mail, Password string
}

//PatchProfile returns the profile of the currently logged in user to the client
func PatchProfile(w http.ResponseWriter, r *http.Request) {
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
		dev.ReportError(err, w, "could not open avatar File"+err.Error())
		return
	}
	io.Copy(file, tempFile)
	w.WriteHeader(200)
}

//RemoveProfilePicture delets the profile picture of the current user
func RemoveProfilePicture(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {
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
	if len(hashedPassword) > 0 {
		if db.UpdatePassword(user.ID, hashedPassword) != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}
	}

	json.NewEncoder(w).Encode(user)
}
