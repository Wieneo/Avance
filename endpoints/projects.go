package endpoints

import (
	"encoding/json"
	"net/http"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
)

//GetProjects returns all projects visible to that user
func GetProjects(w http.ResponseWriter, r *http.Request) {
	if user, err := GetUser(r, w); err == nil {
		projects, err := perms.GetVisibleProjects(user)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(w, err.Error())
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(struct {
			Projects []models.Project
		}{
			projects,
		})
	}

}
