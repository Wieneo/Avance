package perms

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetVisibleProjects returns all projects visible to the user
func GetVisibleProjects(User models.User) ([]models.Project, error) {
	Perms, err := CombinePermissions(User)
	if err != nil {
		return make([]models.Project, 0), err
	}

	if Perms.Admin {
		projects, err := db.GetAllProjects()
		if err != nil {
			return make([]models.Project, 0), err
		}

		return projects, nil
	}

	VisibleProjects := make([]models.Project, 0)
	for _, k := range Perms.AccessTo.Projects {
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
