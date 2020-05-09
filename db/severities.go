package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetSeverities returns all severities from the database relating to the given project
func GetSeverities(Project models.Project, ShowDisabled bool) ([]models.Severity, error) {
	severities := make([]models.Severity, 0)

	/*We double compare here to only use one query
	If ShowDisabled is true it says "Enabled" can to be true and "Enabled" can be false (ShowDisabled Negative)
	If ShowDisabled is false it says "Enabled" can to be true and "Enabled" can be true (ShowDisabled Negative)
	*/
	rows, err := Connection.Query(`SELECT "ID", "Enabled", "Name", "DisplayColor", "Priority" FROM "Severities" WHERE "Enabled" = true OR "Enabled" = $1 AND "Project" = $2`, !ShowDisabled, Project.ID)
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
