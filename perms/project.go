package perms

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"

	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//CheckAccessToProject is used to check if a user has access to a project
func CheckAccessToProject(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check if a project API is called
		if matched, _ := regexp.MatchString("/api/v1/project/[0-9]*/", r.URL.String()); matched {
			projectid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

			project, found, err := db.GetProject(projectid)
			if !found {
				w.WriteHeader(404)
				dev.ReportUserError(w, templates.ProjectNotFound)
				return
			}

			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			user, err := utils.GetUser(r, w)
			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			allperms, perms, err := GetPermissionsToProject(user, project)
			if err != nil {

				utils.ReportInternalErrorToUser(err, w)
				return
			}

			if perms.CanSee || allperms.Admin {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(403)
				dev.ReportUserError(w, templates.ProjectNoPerms)
				return
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

//GetPermissionsToProject returns the ProjectPermission struct regarding the given user and project
func GetPermissionsToProject(User models.User, Project models.Project) (models.Permissions, *models.ProjectPermission, error) {
	perms, err := CombinePermissions(User)
	if err != nil {
		return models.Permissions{}, &models.ProjectPermission{}, err
	}

	if found, pp := containsProject(Project, perms.AccessTo.Projects); found {
		return perms, pp, nil
	}

	return perms, &models.ProjectPermission{}, nil
}
