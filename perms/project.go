package perms

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//CheckAccessToProject is used to check if a user has access to a project
func CheckAccessToProject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if a project API is called
		if matched, _ := regexp.MatchString("/api/v1/project/[0-9]*/", r.URL.String()); matched {
			projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

			project, err := db.GetProject(projectid)
			if err != nil {
				w.WriteHeader(404)
				dev.ReportUserError(w, "Project not found")
				return
			}

			user, err := utils.GetUser(r, w)
			if err != nil {
				w.WriteHeader(500)
				dev.ReportError(err, w, err.Error())
				return
			}

			perms, err := getPermissionsToProject(user, project)
			if err != nil {
				w.WriteHeader(500)
				dev.ReportError(err, w, err.Error())
				return
			}

			if perms.CanSee {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(401)
				dev.ReportUserError(w, "You dont have access to that project")
				return
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func getPermissionsToProject(User models.User, Project models.Project) (*models.ProjectPermission, error) {
	perms, err := combinePermissions(User)
	if err != nil {
		return &models.ProjectPermission{}, err
	}

	if found, pp := containsProject(Project, perms.AccessTo.Projects); found {
		return pp, nil
	}

	return &models.ProjectPermission{}, nil
}
