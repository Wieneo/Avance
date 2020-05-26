package endpoints

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"

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
	println(handler.Filename)
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
