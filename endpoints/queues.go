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
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
)

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

		json.NewEncoder(w).Encode(queues)
	}
}

type queueWebRequest struct {
	Name string
}

//CreateQueue creates a queue
func CreateQueue(w http.ResponseWriter, r *http.Request) {
	var req queueWebRequest

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

	if utils.IsEmpty(req.Name) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Name can't be empty")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Project not found")
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if perms.CanCreateQueues || allperms.Admin {
		queues, err := db.QueuesInProject(project)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		found := false
		for _, k := range queues {
			if strings.ToLower(k.Name) == strings.ToLower(req.Name) {
				found = true
				break
			}
		}

		if found {
			w.WriteHeader(409)
			dev.ReportUserError(w, "A queue with that name already exists")
			return
		}

		id, err := db.CreateQueue(req.Name, projectid)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		queue, found, err := db.GetQueue(projectid, id)
		//If queue isnt found here something went horribly wrong -> ReportError
		if err != nil || !found {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(queue)

	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to create queues in this project")
		return
	}
}

//PatchQueue updates a queue
func PatchQueue(w http.ResponseWriter, r *http.Request) {
	var req queueWebRequest
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Project not found")
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if perms.CanModifyQueues || allperms.Admin {
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

		queue, found, err := db.GetQueue(projectid, queueid)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Queue not found")
			return
		}

		somethingChanged := false

		//Now check if value was specified
		if !utils.IsEmpty(req.Name) && queue.Name != req.Name {
			queue.Name = req.Name
			somethingChanged = true
		}

		if !somethingChanged {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Nothing changed")
			return
		}

		if err := db.PatchQueue(queue); err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(queue)

	} else {
		w.WriteHeader(401)
		dev.ReportUserError(w, "You are not allowed to patch queues in this project")
		return
	}
}

//DeleteQueue deletes a queue
func DeleteQueue(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	project, found, err := db.GetProject(projectid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Project not found")
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if perms.CanRemoveQueues || allperms.Admin {
		_, found, err := db.GetQueue(projectid, queueid)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Queue not found")
			return
		}

		err = db.RemoveQueue(projectid, queueid)
		if err != nil {
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(struct {
			Queue string
		}{
			fmt.Sprintf("Queue %d deleted", queueid),
		})
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to delete queues in this project")
		return
	}
}
