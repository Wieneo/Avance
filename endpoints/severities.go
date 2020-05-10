package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
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
	project, err := db.GetProject(projectid)

	if err != nil {
		w.WriteHeader(404)
		dev.ReportUserError(w, err.Error())
		return
	}

	severities, err := db.GetSeverities(project, showDisabled)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(struct {
		Severities []models.Severity
	}{
		severities,
	})
}

//CreateSeverity creates a severity
func CreateSeverity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Enabled      bool
		Name         string
		DisplayColor string
		Priority     int
	}

	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

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

	if len(req.Name) == 0 || len(req.DisplayColor) == 0 {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Name / DisplayColor can't be empty")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	project, err := db.GetProject(projectid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if perms.CanCreateSeverities {
		severities, err := db.GetSeverities(project, true)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
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
			w.WriteHeader(400)
			dev.ReportUserError(w, "A severity with that name already exists")
			return
		}

		id, err := db.CreateSeverity(req.Enabled, req.Name, req.DisplayColor, req.Priority, projectid)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(struct {
			Severity string
		}{
			fmt.Sprintf("Severity %d created", id),
		})

	} else {
		w.WriteHeader(401)
		dev.ReportUserError(w, "You are not allowed to create severities in this project")
		return
	}
}
