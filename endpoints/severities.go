package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

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

	severities, err := db.GetAllSeverities(showDisabled)
	if err != nil {
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, "Internal Error:"+err.Error())
			return
		}
	}

	json.NewEncoder(w).Encode(struct {
		Severities []models.Severity
	}{
		severities,
	})
}
