package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
			dev.ReportError(err, w, err.Error())
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

//GetProjectQueues returns all queues a user has access to from one project
func GetProjectQueues(w http.ResponseWriter, r *http.Request) {
	if user, err := GetUser(r, w); err == nil {
		//strconv should never throw error because http router expression specifies that only /api/v1/project/[0-9]{*}/queues should be sent here
		projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
		queues, err := perms.GetVisibleQueuesFromProject(user, projectid)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(struct {
			Queues []models.Queue
		}{
			queues,
		})
	}
}
