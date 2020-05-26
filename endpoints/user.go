package endpoints

import (
	"encoding/json"
	"fmt"
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

//GetProfilePicture returns the Profile Picture of the UserID
func GetProfilePicture(w http.ResponseWriter, r *http.Request) {
	userID, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

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
