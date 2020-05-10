package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetSeverities returns all severities from the database relating to the given project
func GetSeverities(Project int64, ShowDisabled bool) ([]models.Severity, error) {
	severities := make([]models.Severity, 0)

	/*We double compare here to only use one query
	If ShowDisabled is true it says "Enabled" can to be true and "Enabled" can be false (ShowDisabled Negative)
	If ShowDisabled is false it says "Enabled" can to be true and "Enabled" can be true (ShowDisabled Negative)
	*/
	rows, err := Connection.Query(`SELECT "ID", "Enabled", "Name", "DisplayColor", "Priority" FROM "Severities" WHERE ("Enabled" = true OR "Enabled" = $1) AND "Project" = $2`, !ShowDisabled, Project)
	if err != nil {
		return severities, err
	}

	for rows.Next() {
		var singleSeverity models.Severity
		rows.Scan(&singleSeverity.ID, &singleSeverity.Enabled, &singleSeverity.Name, &singleSeverity.DisplayColor, &singleSeverity.Priority)
		severities = append(severities, singleSeverity)
	}

	return severities, nil
}

//GetSeverity returns the severity struct to the given severityid
//Throws an error if no severity is found
func GetSeverity(Project int64, Severity int) (models.Severity, bool, error) {
	//I know that project isn't really needed as severity ids are unique anyway
	//Its just a safety measure ;)
	var severity models.Severity
	err := Connection.QueryRow(`SELECT "ID", "Enabled", "Name", "DisplayColor", "Priority" FROM "Severities" WHERE "ID" = $1 AND "Project" = $2`, Severity, Project).Scan(&severity.ID, &severity.Enabled, &severity.Name, &severity.DisplayColor, &severity.Priority)

	if err == nil {
		return severity, true, err
	}

	if err.Error() == "sql: no rows in result set" {
		return severity, false, nil
	}

	return severity, false, err
}

//CreateSeverity creates a severity in the database
func CreateSeverity(Enabled bool, Name, DisplayColor string, Priority int, Project int64) (int, error) {
	var newID int
	err := Connection.QueryRow(`INSERT INTO "Severities" ("Enabled", "Name", "DisplayColor", "Priority", "Project") VALUES ($1, $2, $3, $4, $5) RETURNING "ID"`, Enabled, Name, DisplayColor, Priority, Project).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

//RemoveSeverity removes a severity
func RemoveSeverity(Project int64, Severity int) error {
	//I know that project isn't really needed as severity ids are unique anyway
	//Its just a safety measure ;)
	_, err := Connection.Exec(`DELETE FROM "Severities" WHERE "ID" = $1 AND "Project" = $2`, Severity, Project)
	return err
}
