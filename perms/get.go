package perms

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetVisibleProjects returns all projects visible to the user
func GetVisibleProjects(User models.User) ([]models.Project, error) {
	VisibleProjects := make([]models.Project, 0)
	for _, k := range User.Permissions.AccessTo.Projects {
		if k.CanSee {
			project, err := db.GetProject(k.ProjectID)
			if err != nil {
				return make([]models.Project, 0), err
			}
			VisibleProjects = append(VisibleProjects, project)
		}
	}
	return VisibleProjects, nil
}
