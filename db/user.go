package db

import (
	"encoding/json"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"golang.org/x/crypto/bcrypt"
)

//GetUser returns the user struct from the database
func GetUser(UserID int64) (models.User, error) {
	var Requested models.User
	var RawPermissions string
	err := Connection.QueryRow(`SELECT "ID","Username","Mail", "Permissions", "Firstname", "Lastname" FROM "Users" WHERE "ID" = $1`, UserID).Scan(&Requested.ID, &Requested.Username, &Requested.Mail, &RawPermissions, &Requested.Firstname, &Requested.Lastname)
	if err != nil {
		return Requested, err
	}

	if err := json.Unmarshal([]byte(RawPermissions), &Requested.Permissions); err != nil {
		return Requested, err
	}

	return Requested, nil
}

//GetGroups returns all groups from a user
func GetGroups(User models.User) ([]models.Group, error) {
	Groups := make([]models.Group, 0)
	rows, err := Connection.Query(`SELECT "g"."ID", "g"."Name", "g"."Permissions" FROM "map_User_Group" AS "m" INNER JOIN "Groups" AS "g" ON "g"."ID" = "m"."GroupID" WHERE "m"."UserID" = $1`, User.ID)
	if err != nil {
		return make([]models.Group, 0), err
	}

	for rows.Next() {
		var SingleGroup models.Group
		var RAWJson string
		rows.Scan(&SingleGroup.ID, &SingleGroup.Name, &RAWJson)
		err = json.Unmarshal([]byte(RAWJson), &SingleGroup.Permissions)
		if err != nil {
			return make([]models.Group, 0), err
		}
		Groups = append(Groups, SingleGroup)
	}

	rows.Close()
	return Groups, nil
}

//CreateUser creates a user in the database
func CreateUser(User models.User, Password string) (int64, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(Password), 14)
	if err != nil {
		return 0, err
	}

	permsJSON, err := json.Marshal(User.Permissions)
	if err != nil {
		return 0, err
	}

	var newID int64
	err = Connection.QueryRow(`INSERT INTO "Users" ("Username", "Password", "Mail", "Permissions", "Firstname", "Lastname") VALUES ($1, $2, $3, $4, $5, $6) RETURNING "ID"`, User.Username, string(hash), User.Mail, permsJSON, User.Firstname, User.Lastname).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}
