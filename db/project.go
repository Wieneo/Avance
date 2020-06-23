package db

import (
	"errors"
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetProject returns the project struct to a given projectid
func GetProject(ProjectID int64) (models.Project, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting project %d", ProjectID))
	var Requested models.Project
	rows, err := Connection.Query(`SELECT "ID", "Name", "Description" FROM "Projects" WHERE "ID" = $1`, ProjectID)
	if err != nil {
		return Requested, false, err
	}

	if !rows.Next() {
		dev.LogDebug(fmt.Sprintf("[DB] Projcet %d was not found", ProjectID))
		return Requested, false, errors.New("Project not found")
	}

	rows.Scan(&Requested.ID, &Requested.Name, &Requested.Description)
	rows.Close()

	dev.LogDebug(fmt.Sprintf("[DB] Got Project %d: Name: %s, Description: %s", ProjectID, Requested.Name, Requested.Description))
	return Requested, true, nil
}

//GetAllProjects returns all projects
func GetAllProjects() ([]models.Project, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting ALL Projects"))
	AllProjects := make([]models.Project, 0)
	rows, err := Connection.Query(`SELECT "ID", "Name", "Description" FROM "Projects"`)

	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while gettings all projects -> Returning empty array (%s)", err.Error()))
		return make([]models.Project, 0), err
	}

	for rows.Next() {
		var SingleProject models.Project
		rows.Scan(&SingleProject.ID, &SingleProject.Name, &SingleProject.Description)
		AllProjects = append(AllProjects, SingleProject)
	}

	rows.Close()

	dev.LogDebug(fmt.Sprintf("[DB] Got %d projects", len(AllProjects)))

	return AllProjects, nil
}

//CreateProject creates a project in the database
func CreateProject(Name, Description string) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Creating project with values: Name: %s, Description: %s", Name, Description))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Projects" ("Name", "Description") VALUES ($1, $2) RETURNING "ID"`, Name, Description).Scan(&newID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Project creation failed -> Returning 0 Value (%s)", err.Error()))
		return 0, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Created project with id %d", newID))

	return newID, nil
}

//PatchProject updates the project in the database
func PatchProject(Project models.Project) error {
	dev.LogDebug(fmt.Sprintf("[DB] Patching project %d", Project.ID))
	_, err := Connection.Exec(`UPDATE "Projects" SET "Name" = $1, "Description" = $2 WHERE "ID" = $3`, Project.Name, Project.Description, Project.ID)
	dev.LogDebug(fmt.Sprintf("[DB] Project %d patched", Project.ID))
	return err
}
