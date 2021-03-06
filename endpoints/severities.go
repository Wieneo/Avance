package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
)

//GetSeverities returns all severities
func GetSeverities(w http.ResponseWriter, r *http.Request) {
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
	project, found, err := db.GetProject(projectid)

	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	severities, err := db.GetSeverities(project.ID, showDisabled)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(severities)
}

type severityWebRequest struct {
	Enabled      bool
	Name         string
	DisplayColor string
	Priority     int
}

//CreateSeverity creates a severity
func CreateSeverity(w http.ResponseWriter, r *http.Request) {
	var req severityWebRequest

	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

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

	if utils.IsEmpty(req.Name) || utils.IsEmpty(req.DisplayColor) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Name / DisplayColor can't be empty")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if perms.CanCreateSeverities || allperms.Admin {
		severities, err := db.GetSeverities(project.ID, true)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		found := false
		for _, k := range severities {
			if strings.ToLower(k.Name) == strings.ToLower(req.Name) {
				found = true
				break
			}
		}

		if found {
			w.WriteHeader(409)
			dev.ReportUserError(w, "A severity with that name already exists")
			return
		}

		id, err := db.CreateSeverity(req.Enabled, req.Name, req.DisplayColor, req.Priority, projectid)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		sev, found, err := db.GetSeverity(projectid, id)
		//If severity isnt found here something went horribly wrong -> ReportError
		if err != nil || !found {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(sev)

	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to create severities in this project")
		return
	}
}

//PatchSeverity updates a severity
func PatchSeverity(w http.ResponseWriter, r *http.Request) {
	var req severityWebRequest
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	severityid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if perms.CanModifySeverities || allperms.Admin {
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

		severity, found, err := db.GetSeverity(projectid, severityid)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, templates.SeverityNotFound)
			return
		}

		somethingChanged := false

		//Now check if value was specified
		if !utils.IsEmpty(req.Name) && severity.Name != req.Name {
			severity.Name = req.Name
			somethingChanged = true
		}

		if !utils.IsEmpty(req.DisplayColor) && severity.DisplayColor != req.DisplayColor {
			severity.DisplayColor = req.DisplayColor
			somethingChanged = true
		}

		//Check for occurence in string (request body) as we cant differenciate if the value was specified or not
		if strings.Contains(string(rawBytes), "Enabled") {
			severity.Enabled = req.Enabled
			somethingChanged = true
		}

		if strings.Contains(string(rawBytes), "Priority") {
			severity.Priority = req.Priority
			somethingChanged = true
		}

		if !somethingChanged {
			w.WriteHeader(406)
			dev.ReportUserError(w, templates.NothingChanged)
			return
		}

		if err := db.PatchSeverity(severity); err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(severity)

	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to patch severities in this project")
		return
	}
}

//DeleteSeverity deletes a severity
func DeleteSeverity(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	severityid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.ProjectNotFound)
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if perms.CanRemoveSeverities || allperms.Admin {
		_, found, err := db.GetSeverity(projectid, severityid)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, templates.SeverityNotFound)
			return
		}

		err = db.RemoveSeverity(projectid, severityid)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(struct {
			Severity string
		}{
			fmt.Sprintf("Severity %d deleted", severityid),
		})
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to delete severities in this project")
		return
	}
}
