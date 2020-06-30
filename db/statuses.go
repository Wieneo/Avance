package db

import (
	"database/sql"
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//SearchStatus searches for a status and returns the ID, if a status was found and maybe an error
func SearchStatus(Project int64, Name string) (int64, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Searching for status '%s' in project %d", Name, Project))
	var ID int64

	//Ignoring casing
	err := Connection.QueryRow(`SELECT "ID" FROM "Statuses" WHERE UPPER("Name") = UPPER($1) AND "Project" = $2`, Name, Project).Scan(&ID)
	if err != nil {
		if err == sql.ErrNoRows {
			dev.LogDebug(fmt.Sprintf("[DB] Found no status with name '%s' in project %d", Name, Project))
			return ID, false, nil
		}

		return ID, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Found status %d in project %d", ID, Project))
	return ID, true, nil
}

//GetStatuses returns all statuses from the database relating to the given project
func GetStatuses(Project int64, ShowDisabled bool) ([]models.Status, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting ALL statuses from project %d (Showing disabled: %t)", Project, ShowDisabled))

	statuses := make([]models.Status, 0)

	/*We double compare here to only use one query
	If ShowDisabled is true it says "Enabled" can to be true and "Enabled" can be false (ShowDisabled Negative)
	If ShowDisabled is false it says "Enabled" can to be true and "Enabled" can be true (ShowDisabled Negative)
	*/
	rows, err := Connection.Query(`SELECT "ID", "Enabled", "Name", "DisplayColor", "TicketsVisible" FROM "Statuses" WHERE ("Enabled" = true OR "Enabled" = $1) AND "Project" = $2`, !ShowDisabled, Project)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while getting all statuses from project %d: %s", Project, err.Error()))
		return statuses, err
	}

	for rows.Next() {
		var singleStatus models.Status
		rows.Scan(&singleStatus.ID, &singleStatus.Enabled, &singleStatus.Name, &singleStatus.DisplayColor, &singleStatus.TicketsVisible)
		statuses = append(statuses, singleStatus)
	}

	rows.Close()

	dev.LogDebug(fmt.Sprintf("[DB] Got %d statuses from project %d", len(statuses), Project))
	return statuses, nil
}

//GetStatus returns the status struct to the given statusid
//Throws an error if no status is found
func GetStatus(Project, Status int64) (models.Status, bool, error) {
	//I know that project isn't really needed as status ids are unique anyway
	//Its just a safety measure ;)
	dev.LogDebug(fmt.Sprintf("[DB] Getting status %d in project %d", Status, Project))
	var status models.Status
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "TicketsVisible" FROM "Statuses" WHERE "ID" = $1 AND "Project" = $2`, Status, Project).Scan(&status.ID, &status.Enabled, &status.Name, &status.DisplayColor, &status.TicketsVisible)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Got status %s (ID: %d) from project %d", status.Name, status.ID, Project))
		return status, true, err
	}

	if err == sql.ErrNoRows {
		dev.LogDebug(fmt.Sprintf("[DB] Status %d not found in project %d", Status, Project))
		return status, false, nil
	}

	dev.LogDebug(fmt.Sprintf("[DB] Error happened while getting status %d from project %d: %s", Status, Project, err.Error()))
	return status, false, err
}

//GetStatusUNSAFE is a copy of GetStatus but without the ProjectID given
func GetStatusUNSAFE(Status int64) (models.Status, bool, error) {
	//I know that project isn't really needed as status ids are unique anyway
	//Its just a safety measure ;)
	dev.LogDebug(fmt.Sprintf("[DB] Getting status %d in UNSAFE mode", Status))
	var status models.Status
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "TicketsVisible" FROM "Statuses" WHERE "ID" = $1`, Status).Scan(&status.ID, &status.Enabled, &status.Name, &status.DisplayColor, &status.TicketsVisible)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Got status %s (ID: %d)", status.Name, status.ID))
		return status, true, err
	}

	if err == sql.ErrNoRows {
		dev.LogDebug(fmt.Sprintf("[DB] Status %d not found", Status))
		return status, false, nil
	}

	dev.LogDebug(fmt.Sprintf("[DB] Error happened while getting status %d: %s", Status, err.Error()))
	return status, false, err
}

//CreateStatus creates a status in the database
func CreateStatus(Enabled bool, Name, DisplayColor string, TicketsVisible bool, Project int64) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Creating status with values: Name: %s, DisplayColor: %s, TicketVisible: %t in project %d", Name, DisplayColor, TicketsVisible, Project))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Statuses" ("Enabled", "Name", "DisplayColor", "TicketsVisible", "Project") VALUES ($1, $2, $3, $4, $5) RETURNING "ID"`, Enabled, Name, DisplayColor, TicketsVisible, Project).Scan(&newID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while creating status in project %d: %s", Project, err.Error()))
		return 0, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Created status (ID: %d) in project %d", newID, Project))
	return newID, nil
}

//PatchStatus updates the give status
func PatchStatus(Status models.Status) error {
	dev.LogDebug(fmt.Sprintf("[DB] Patching status %d", Status.ID))
	_, err := Connection.Exec(`UPDATE "Statuses" SET "Enabled" = $1, "Name" = $2, "DisplayColor" = $3, "TicketsVisible" = $4 WHERE "ID" = $5`, Status.Enabled, Status.Name, Status.DisplayColor, Status.TicketsVisible, Status.ID)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Patched status %d", Status.ID))
	}
	return err
}

//RemoveStatus removes a status
func RemoveStatus(Project, Status int64) error {
	//I know that project isn't really needed as status ids are unique anyway
	//Its just a safety measure ;)
	dev.LogDebug(fmt.Sprintf("[DB] Removing status %d from project %d", Status, Project))
	_, err := Connection.Exec(`DELETE FROM "Statuses" WHERE "ID" = $1 AND "Project" = $2`, Status, Project)
	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Removed status %d from project %d", Status, Project))
	}
	return err
}
