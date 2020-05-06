package db

import (
	"errors"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetProject returns the project struct to a given projectid
func GetProject(ProjectID int64) (models.Project, error) {
	var Requested models.Project
	rows, err := Connection.Query(`SELECT * FROM "Projects" WHERE "ID" = $1`, ProjectID)
	if err != nil {
		return Requested, err
	}

	if !rows.Next() {
		return Requested, errors.New("Project not found")
	}

	rows.Scan(&Requested.ID, &Requested.Name, &Requested.Description)
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

//QueuesInProject returns all queues from a project
func QueuesInProject(Project models.Project) ([]models.Queue, error) {
	Queues := make([]models.Queue, 0)
	rows, err := Connection.Query(`SELECT "ID", "Name" FROM "Queue" WHERE "Project" = $1`, Project.ID)
	if err != nil {
		return make([]models.Queue, 0), err
	}

	for rows.Next() {
		var SingleQueue models.Queue
		rows.Scan(&SingleQueue.ID, &SingleQueue.Name)
		Queues = append(Queues, SingleQueue)
	}

	rows.Close()

	return Queues, nil
}
