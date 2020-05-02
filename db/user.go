package db

import (
	"encoding/json"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetUser returns the user struct from the database
func GetUser(UserID int) (models.User, error) {
	var Requested models.User
	var RawPermissions string
	err := Connection.QueryRow(`SELECT "ID","Username","Mail", "Permissions" FROM "Users" WHERE "ID" = $1`, UserID).Scan(&Requested.ID, &Requested.Username, &Requested.Mail, &RawPermissions)
	if err != nil {
		return Requested, err
	}

	if err := json.Unmarshal([]byte(RawPermissions), &Requested.Permissions); err != nil {
		return Requested, err
	}

	return Requested, nil
}
