package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
