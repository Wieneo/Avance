package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
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
		w.WriteHeader(400)
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
