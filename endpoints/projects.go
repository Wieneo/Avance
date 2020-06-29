package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//GetProjects returns all projects visible to that user
func GetProjects(w http.ResponseWriter, r *http.Request) {
	if user, err := utils.GetUser(r, w); err == nil {
		projects, err := perms.GetVisibleProjects(user)
		if err != nil {
			utils.ReportErrorToUser(err, w)
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(projects)
	}

}

//GetSingleProject returns all projects visible to that user
func GetSingleProject(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {
		utils.ReportErrorToUser(err, w)
		return
	}

	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	project, found, err := db.GetProject(projectid)

	if err != nil {
		utils.ReportErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Project not found")
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {
		utils.ReportErrorToUser(err, w)
		return
	}

	if perms.CanSee || allperms.Admin {
		json.NewEncoder(w).Encode(project)
	} else {
		w.WriteHeader(401)
		dev.ReportUserError(w, "You don't have access to that project")
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
		utils.ReportErrorToUser(err, w)
		return
	}

	err = json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Request is malformed: "+err.Error())
		return
	}

	if utils.IsEmpty(req.Name) {
		w.WriteHeader(400)
		dev.ReportUserError(w, "Name can't be empty")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		utils.ReportErrorToUser(err, w)
		return
	}

	//Projects only can be created by administrators
	if admin, err := perms.IsAdmin(user); err == nil {
		if admin {
			projects, err := db.GetAllProjects()
			if err != nil {
				utils.ReportErrorToUser(err, w)
				return
			}

			found := false
			for _, k := range projects {
				if strings.ToLower(k.Name) == strings.ToLower(req.Name) {
					found = true
				}
			}

			if found {
				w.WriteHeader(403)
				dev.ReportUserError(w, "Project with that name already exists")
				return
			}

			id, err := db.CreateProject(req.Name, req.Description)
			if err != nil {
				utils.ReportErrorToUser(err, w)
				return
			}

			project, found, err := db.GetProject(id)
			//If the project is not found something went horribly wrong -> ReportError here
			if err != nil || !found {
				utils.ReportErrorToUser(err, w)
				return
			}

			json.NewEncoder(w).Encode(project)

		} else {
			w.WriteHeader(401)
			dev.ReportUserError(w, "You are not allowed to create projects")
			return
		}
	} else {
		utils.ReportErrorToUser(err, w)
		return
	}
}

//ChangeProject updates the given project
func ChangeProject(w http.ResponseWriter, r *http.Request) {
	var req projectWebRequest
	projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

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
		dev.ReportUserError(w, "Project not found")
		return
	}

	allperms, perms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {
		utils.ReportErrorToUser(err, w)
		return
	}

	if perms.CanModify || allperms.Admin {
		//Parse JSON
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

		somethingChanged := false

		//Data in req variable
		//Check for occurence in string (request body) as we cant differenciate if the value was specified or not
		if strings.Contains(string(rawBytes), "Name") {
			project.Name = req.Name
			somethingChanged = true
		}

		if strings.Contains(string(rawBytes), "Description") {
			project.Description = req.Description
			somethingChanged = true
		}

		if !somethingChanged {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Nothing changed")
			return
		}

		err = db.PatchProject(project)
		if err != nil {
			utils.ReportErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(project)
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to patch this project")
		return
	}
}
