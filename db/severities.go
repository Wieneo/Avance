package db

import (
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//SearchSeverity searches for a severity and returns the ID, if a severity was found and maybe an error
func SearchSeverity(Project int64, Name string) (int64, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Searching for severity '%s' in project %d", Name, Project))
	var ID int64

	//Ignoring casing
	err := Connection.QueryRow(`SELECT "ID" FROM "Severities" WHERE UPPER("Name") = UPPER($1) AND "Project" = $2`, Name, Project).Scan(&ID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			dev.LogDebug(fmt.Sprintf("[DB] Severity '%s' wasn't found in project %d", Name, Project))
			return ID, false, nil
		}

		return ID, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Query '%s' returned severity %d", Name, ID))
	return ID, true, nil
}

//GetSeverities returns all severities from the database relating to the given project
func GetSeverities(Project int64, ShowDisabled bool) ([]models.Severity, error) {
	if ShowDisabled {
		dev.LogDebug(fmt.Sprintf("[DB] Getting ALL severities in project %d (Showing disabled: true)", Project))
	} else {
		dev.LogDebug(fmt.Sprintf("[DB] Getting ALL severities in project %d (Showing disabled: false)", Project))
	}
	severities := make([]models.Severity, 0)

	/*We double compare here to only use one query
	If ShowDisabled is true it says "Enabled" can to be true and "Enabled" can be false (ShowDisabled Negative)
	If ShowDisabled is false it says "Enabled" can to be true and "Enabled" can be true (ShowDisabled Negative)
	*/
	rows, err := Connection.Query(`SELECT "ID", "Enabled", "Name", "DisplayColor", "Priority" FROM "Severities" WHERE ("Enabled" = true OR "Enabled" = $1) AND "Project" = $2`, !ShowDisabled, Project)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while getting all severities in project %d: %s", Project, err.Error()))
		return severities, err
	}

	for rows.Next() {
		var singleSeverity models.Severity
		rows.Scan(&singleSeverity.ID, &singleSeverity.Enabled, &singleSeverity.Name, &singleSeverity.DisplayColor, &singleSeverity.Priority)
		severities = append(severities, singleSeverity)
	}

	rows.Close()
	dev.LogDebug(fmt.Sprintf("[DB] Got %d severities in project %d", len(severities), Project))
	return severities, nil
}

//GetSeverity returns the severity struct to the given severityid
//Throws an error if no severity is found
func GetSeverity(Project int64, Severity int64) (models.Severity, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting severity %d in project %d", Severity, Project))
	//I know that project isn't really needed as severity ids are unique anyway
	//Its just a safety measure ;)
	var severity models.Severity
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "Priority" FROM "Severities" WHERE "ID" = $1 AND "Project" = $2`, Severity, Project).Scan(&severity.ID, &severity.Enabled, &severity.Name, &severity.DisplayColor, &severity.Priority)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Got severity %d successfully", Severity))
		return severity, true, err
	}

	if err.Error() == "sql: no rows in result set" {
		dev.LogDebug(fmt.Sprintf("[DB] Severity %d wasn't found in project %d", Severity, Project))
		return severity, false, nil
	}

	dev.LogDebug(fmt.Sprintf("[DB] There was an error getting severity %d: %s", Severity, err.Error()))
	return severity, false, err
}

//GetSeverityUNSAFE is a copy of GetSeverity but without the ProjectID needed
func GetSeverityUNSAFE(Severity int64) (models.Severity, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting severity %d in UNSAFE mode", Severity))
	//I know that project isn't really needed as severity ids are unique anyway
	//Its just a safety measure ;)
	var severity models.Severity
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "Priority" FROM "Severities" WHERE "ID" = $1`, Severity).Scan(&severity.ID, &severity.Enabled, &severity.Name, &severity.DisplayColor, &severity.Priority)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Got severity %d successfully", Severity))
		return severity, true, err
	}

	if err.Error() == "sql: no rows in result set" {
		dev.LogDebug(fmt.Sprintf("[DB] Severity %d wasn't found", Severity))
		return severity, false, nil
	}

	dev.LogDebug(fmt.Sprintf("[DB] There was an error getting severity %d: %s", Severity, err.Error()))
	return severity, false, err
}

//CreateSeverity creates a severity in the database
func CreateSeverity(Enabled bool, Name, DisplayColor string, Priority int, Project int64) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Creating severity in project %d with values: Name: %s, DisplayColor: %s, Priority: %d", Project, Name, DisplayColor, Priority))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Severities" ("Enabled", "Name", "DisplayColor", "Priority", "Project") VALUES ($1, $2, $3, $4, $5) RETURNING "ID"`, Enabled, Name, DisplayColor, Priority, Project).Scan(&newID)
	if err != nil {
		return 0, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Created severity in project %d with id %d", Project, newID))
	return newID, nil
}

//PatchSeverity updates the give severity
func PatchSeverity(Severity models.Severity) error {
	dev.LogDebug(fmt.Sprintf("[DB] Patching severity %d", Severity.ID))
	_, err := Connection.Exec(`UPDATE "Severities" SET "Enabled" = $1, "Name" = $2, "DisplayColor" = $3, "Priority" = $4 WHERE "ID" = $5`, Severity.Enabled, Severity.Name, Severity.DisplayColor, Severity.Priority, Severity.ID)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Patched Severity %d", Severity.ID))
	}
	return err
}

//RemoveSeverity removes a severity
func RemoveSeverity(Project int64, Severity int64) error {
	dev.LogDebug(fmt.Sprintf("[DB] Removing severity %d from project %d", Severity, Project))
	//I know that project isn't really needed as severity ids are unique anyway
	//Its just a safety measure ;)
	_, err := Connection.Exec(`DELETE FROM "Severities" WHERE "ID" = $1 AND "Project" = $2`, Severity, Project)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Severity %d in project %d deleted", Severity, Project))
	}
	return err
}
