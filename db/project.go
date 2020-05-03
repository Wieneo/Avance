package db

import "gitlab.gnaucke.dev/tixter/tixter-app/v2/models"

//GetProject returns the project struct to a given projectid
func GetProject(ProjectID int) (models.Project, error) {
	var Requested models.Project
	err := Connection.QueryRow(`SELECT * FROM "Projects" WHERE "ID" = $1`, ProjectID).Scan(&Requested.ID, &Requested.Name)
	if err != nil {
		return Requested, err
	}
	return Requested, nil
}

//GetAllProjects returns all projects
func GetAllProjects() ([]models.Project, error) {
	AllProjects := make([]models.Project, 0)
	rows, err := Connection.Query(`SELECT * FROM "Projects"`)

	if err != nil {
		return make([]models.Project, 0), err
	}

	for rows.Next() {
		var SingleProject models.Project
		rows.Scan(&SingleProject.ID, &SingleProject.Name)
		AllProjects = append(AllProjects, SingleProject)
	}

	rows.Close()

	return AllProjects, nil
}
