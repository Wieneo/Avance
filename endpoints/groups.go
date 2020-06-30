package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

type groupWebRequest struct {
	Name string
}

//GetGroups returns all groups
func GetGroups(w http.ResponseWriter, r *http.Request) {
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
	if !perms.Admin && !perms.CanChangePermissionsGlobal && !perms.CanModifyGroups {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to view all groups")
		return
	}

	groups, err := db.GetALLGroups()
	if err != nil {
		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(groups)
}

//CreateGroup updates profile information
func CreateGroup(w http.ResponseWriter, r *http.Request) {
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

	if !perms.CanCreateGroups && !perms.Admin {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed create groups")
		return
	}

	var req groupWebRequest

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

	if utils.IsEmpty(req.Name) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Name can't be empty")
		return
	}

	groups, err := db.GetALLGroups()
	if err != nil {
		utils.ReportInternalErrorToUser(err, w)
		return
	}

	for _, k := range groups {
		if strings.ToLower(k.Name) == strings.ToLower(req.Name) {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Group with that name already exists")
			return
		}
	}

	var newGroup models.Group = models.Group{
		Name: req.Name,
	}

	newGroup.ID, err = db.CreateGroup(newGroup)
	if err != nil {
		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(newGroup)
}
