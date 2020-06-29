package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//GetStatuses returns all statuses
func GetStatuses(w http.ResponseWriter, r *http.Request) {
	//If the user has access to the project has already been checked
	showDisabled := false
	var err error
	if len(r.URL.Query()["showDisabled"]) > 0 {
		showDisabled, err = strconv.ParseBool(r.URL.Query()["showDisabled"][0])
		if err != nil {
			w.WriteHeader(400)
			dev.ReportUserError(w, "showDisabled Argument is not a boolean")
			return
		}
	}

	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

	statuses, err := db.GetStatuses(projectid, showDisabled)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(statuses)
}

type statusWebRequest struct {
	Enabled        bool
	Name           string
	DisplayColor   string
	TicketsVisible bool
}

//CreateStatus creates a status
func CreateStatus(w http.ResponseWriter, r *http.Request) {
	var req statusWebRequest

	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

	rawBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	err = json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Request is malformed: "+err.Error())
		return
	}

	if utils.IsEmpty(req.Name) || utils.IsEmpty(req.DisplayColor) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Name / DisplayColor can't be empty")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if perms.CanCreateStatuses || allperms.Admin {
		statuses, err := db.GetStatuses(project.ID, true)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		found := false
		for _, k := range statuses {
			if strings.ToLower(k.Name) == strings.ToLower(req.Name) {
				found = true
				break
			}
		}

		if found {
			w.WriteHeader(409)
			dev.ReportUserError(w, "A status with that name already exists")
			return
		}

		id, err := db.CreateStatus(req.Enabled, req.Name, req.DisplayColor, req.TicketsVisible, projectid)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		status, found, err := db.GetStatus(projectid, id)
		if err != nil || !found {

			utils.ReportErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(status)

	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to create statuses in this project")
		return
	}
}

//PatchStatus updates a status
func PatchStatus(w http.ResponseWriter, r *http.Request) {
	var req statusWebRequest
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	statusid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if perms.CanModifyStatuses || allperms.Admin {
		rawBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		err = json.Unmarshal(rawBytes, &req)
		if err != nil {
			w.WriteHeader(400)
			dev.ReportUserError(w, "Request is malformed: "+err.Error())
			return
		}

		status, found, err := db.GetStatus(projectid, statusid)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, templates.StatusNotFound)
			return
		}

		somethingChanged := false

		//Now check if value was specified
		if !utils.IsEmpty(req.Name) && status.Name != req.Name {
			status.Name = req.Name
			somethingChanged = true
		}

		if !utils.IsEmpty(req.DisplayColor) && status.DisplayColor != req.DisplayColor {
			status.DisplayColor = req.DisplayColor
			somethingChanged = true
		}

		//Check for occurence in string (request body) as we cant differenciate if the value was specified or not
		if strings.Contains(string(rawBytes), "Enabled") {
			status.Enabled = req.Enabled
			somethingChanged = true
		}

		if strings.Contains(string(rawBytes), "TicketsVisible") {
			status.TicketsVisible = req.TicketsVisible
			somethingChanged = true
		}

		if !somethingChanged {
			w.WriteHeader(400)
			dev.ReportUserError(w, templates.NothingChanged)
			return
		}

		if err := db.PatchStatus(status); err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(status)

	} else {
		w.WriteHeader(401)
		dev.ReportUserError(w, "You are not allowed to patch statuses in this project")
		return
	}
}

//DeleteStatus deletes a status
func DeleteStatus(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	statusid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if perms.CanRemoveStatuses || allperms.Admin {
		_, found, err := db.GetStatus(projectid, statusid)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, templates.StatusNotFound)
			return
		}

		err = db.RemoveStatus(projectid, statusid)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(struct {
			Status string
		}{
			fmt.Sprintf("Status %d deleted", statusid),
		})
	} else {
		w.WriteHeader(401)
		dev.ReportUserError(w, "You are not allowed to delete statuses in this project")
		return
	}
}
