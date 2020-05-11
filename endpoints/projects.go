package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
)

//GetProjects returns all projects visible to that user
func GetProjects(w http.ResponseWriter, r *http.Request) {
	if user, err := utils.GetUser(r, w); err == nil {
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

type projectWebRequest struct {
	Name        string
	Description string
}

//CreateProject creates a project in the database
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var req projectWebRequest

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

	if len(req.Name) == 0 {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Name can't be empty")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	//Projects only can be created by administrators
	if admin, err := perms.IsAdmin(user); err == nil {
		if admin {
			projects, err := db.GetAllProjects()
			if err != nil {
				w.WriteHeader(500)
				dev.ReportError(err, w, err.Error())
				return
			}

			found := false
			for _, k := range projects {
				if strings.ToLower(k.Name) == strings.ToLower(req.Name) {
					found = true
				}
			}

			if found {
				w.WriteHeader(401)
				dev.ReportUserError(w, "Project with that name already exists")
				return
			}

			id, err := db.CreateProject(req.Name, req.Description)
			if err != nil {
				w.WriteHeader(500)
				dev.ReportError(err, w, err.Error())
				return
			}

			json.NewEncoder(w).Encode(struct {
				Project string
			}{
				fmt.Sprintf("Project %d created", id),
			})

		} else {
			w.WriteHeader(401)
			dev.ReportUserError(w, "You are not allowed to create projects")
			return
		}
	} else {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}
}

//GetProjectQueues returns all queues a user has access to from one project
func GetProjectQueues(w http.ResponseWriter, r *http.Request) {
	if user, err := utils.GetUser(r, w); err == nil {
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
