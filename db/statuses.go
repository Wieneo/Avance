package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//SearchStatus searches for a status and returns the ID, if a status was found and maybe an error
func SearchStatus(Project int64, Name string) (int64, bool, error) {
	var ID int64

	//Ignoring casing
	err := Connection.QueryRow(`SELECT "ID" FROM "Statuses" WHERE UPPER("Name") = UPPER($1) AND "Project" = $2`, Name, Project).Scan(&ID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ID, false, nil
		}

		return ID, true, err
	}

	return ID, true, nil
}

//GetStatuses returns all statuses from the database relating to the given project
func GetStatuses(Project int64, ShowDisabled bool) ([]models.Status, error) {
	statuses := make([]models.Status, 0)

	/*We double compare here to only use one query
	If ShowDisabled is true it says "Enabled" can to be true and "Enabled" can be false (ShowDisabled Negative)
	If ShowDisabled is false it says "Enabled" can to be true and "Enabled" can be true (ShowDisabled Negative)
	*/
	rows, err := Connection.Query(`SELECT "ID", "Enabled", "Name", "DisplayColor", "TicketsVisible" FROM "Statuses" WHERE ("Enabled" = true OR "Enabled" = $1) AND "Project" = $2`, !ShowDisabled, Project)
	if err != nil {
		return statuses, err
	}

	for rows.Next() {
		var singleStatus models.Status
		rows.Scan(&singleStatus.ID, &singleStatus.Enabled, &singleStatus.Name, &singleStatus.DisplayColor, &singleStatus.TicketsVisible)
		statuses = append(statuses, singleStatus)
	}

	rows.Close()
	return statuses, nil
}

//GetStatus returns the status struct to the given statusid
//Throws an error if no status is found
func GetStatus(Project int64, Status int64) (models.Status, bool, error) {
	//I know that project isn't really needed as status ids are unique anyway
	//Its just a safety measure ;)
	var status models.Status
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "TicketsVisible" FROM "Statuses" WHERE "ID" = $1 AND "Project" = $2`, Status, Project).Scan(&status.ID, &status.Enabled, &status.Name, &status.DisplayColor, &status.TicketsVisible)

	if err == nil {
		return status, true, err
	}

	if err.Error() == "sql: no rows in result set" {
		return status, false, nil
	}

	return status, false, err
}

//GetStatusUNSAFE is a copy of GetStatus but without the ProjectID given
func GetStatusUNSAFE(Status int64) (models.Status, bool, error) {
	//I know that project isn't really needed as status ids are unique anyway
	//Its just a safety measure ;)
	var status models.Status
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "TicketsVisible" FROM "Statuses" WHERE "ID" = $1`, Status).Scan(&status.ID, &status.Enabled, &status.Name, &status.DisplayColor, &status.TicketsVisible)

	if err == nil {
		return status, true, err
	}

	if err.Error() == "sql: no rows in result set" {
		return status, false, nil
	}

	return status, false, err
}

//CreateStatus creates a status in the database
func CreateStatus(Enabled bool, Name, DisplayColor string, TicketsVisible bool, Project int64) (int64, error) {
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Statuses" ("Enabled", "Name", "DisplayColor", "TicketsVisible", "Project") VALUES ($1, $2, $3, $4, $5) RETURNING "ID"`, Enabled, Name, DisplayColor, TicketsVisible, Project).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

//PatchStatus updates the give status
func PatchStatus(Status models.Status) error {
	_, err := Connection.Exec(`UPDATE "Statuses" SET "Enabled" = $1, "Name" = $2, "DisplayColor" = $3, "TicketsVisible" = $4 WHERE "ID" = $5`, Status.Enabled, Status.Name, Status.DisplayColor, Status.TicketsVisible, Status.ID)
	return err
}

//RemoveStatus removes a status
func RemoveStatus(Project int64, Status int64) error {
	//I know that project isn't really needed as status ids are unique anyway
	//Its just a safety measure ;)
	_, err := Connection.Exec(`DELETE FROM "Statuses" WHERE "ID" = $1 AND "Project" = $2`, Status, Project)
	return err
}
