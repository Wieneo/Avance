package db

import (
	"encoding/json"
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetALLGroups returns all groups from the database. This should be used with caution as it can cause many cpu cycles
func GetALLGroups() ([]models.Group, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting ALL Groups"))
	groups := make([]models.Group, 0)
	rows, err := Connection.Query(`SELECT "Name" FROM "Groups"`)
	if err != nil {
		dev.LogError(err, "Error occured while getting group: "+err.Error())
		return make([]models.Group, 0), err
	}

	for rows.Next() {
		var singleUser models.Group
		rows.Scan(&singleUser.Name)
		groups = append(groups, singleUser)
	}

	dev.LogDebug(fmt.Sprintf("[DB] Got %d groups", len(groups)))

	return groups, nil
}

//CreateGroup creates a group in the database
func CreateGroup(Group models.Group) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Creating group '%s'", Group.Name))
	permsJSON, err := json.Marshal(Group.Permissions)
	if err != nil {
		return 0, err
	}

	var newID int64
	err = Connection.QueryRow(`INSERT INTO "Groups" ("Name", "Permissions") VALUES ($1, $2) RETURNING "ID"`, Group.Name, permsJSON).Scan(&newID)
	if err != nil {
		return 0, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Created group %d", newID))

	return newID, nil
}
