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
