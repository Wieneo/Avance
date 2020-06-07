package db

import (
	"encoding/json"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"golang.org/x/crypto/bcrypt"
)

//SearchUser searches for a user and returns the ID, if a user was found and maybe an error
func SearchUser(Name string) (int64, bool, error) {
	var ID int64

	//Ignoring casing
	err := Connection.QueryRow(`SELECT "ID" FROM "Users" WHERE UPPER("Username") = UPPER($1) AND "Active" = true`, Name).Scan(&ID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return ID, false, nil
		}

		return ID, true, err
	}

	return ID, true, nil
}

//GetUser returns the user struct from the database
func GetUser(UserID int64) (models.User, bool, error) {
	return getUser(UserID, true)
}

//DumbGetUser returns the user struct from the database ignoring the "Active" field
func DumbGetUser(UserID int64) (models.User, bool, error) {
	return getUser(UserID, false)
}

func getUser(UserID int64, RespectActive bool) (models.User, bool, error) {
	var Requested models.User
	var RawPermissions string
	err := Connection.QueryRow(`SELECT "ID","Username","Mail", "Permissions", "Firstname", "Lastname" FROM "Users" WHERE "ID" = $1 AND ("Active" = true OR "Active" = $2)`, UserID, RespectActive).Scan(&Requested.ID, &Requested.Username, &Requested.Mail, &RawPermissions, &Requested.Firstname, &Requested.Lastname)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return Requested, false, nil
		}

		return Requested, true, err
	}

	if err := json.Unmarshal([]byte(RawPermissions), &Requested.Permissions); err != nil {
		return Requested, true, err
	}

	//Fix that AccessTo arrays are never null / nil
	if Requested.Permissions.AccessTo.Projects == nil {
		Requested.Permissions.AccessTo.Projects = make([]models.ProjectPermission, 0)
	}

	if Requested.Permissions.AccessTo.Queues == nil {
		Requested.Permissions.AccessTo.Queues = make([]models.QueuePermission, 0)
	}
	////////////////////////////////////////////////

	return Requested, true, nil
}

//GetSettings populates the Settings struct in the User struct
func GetSettings(User *models.User) error {
	var RawSettings string
	err := Connection.QueryRow(`SELECT "Settings" FROM "Users" WHERE "ID" = $1`, User.ID).Scan(&RawSettings)

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(RawSettings), &User.Settings); err != nil {
		return err
	}

	//Make them human readable
	for _, k := range User.Settings.EnabledNotificationChannels {
		User.Settings.EnabledNotificationChannelsReadable = append(User.Settings.EnabledNotificationChannelsReadable, k.String())
	}

	return nil
}

//GetALLUsers returns all users from the database. This should be used with caution as it can cause many cpu cycles
func GetALLUsers() ([]models.User, error) {
	users := make([]models.User, 0)
	rows, err := Connection.Query(`SELECT "ID" FROM "Users" WHERE "Active" = true`)
	if err != nil {
		dev.LogError(err, "Error occured while getting users: "+err.Error())
		return make([]models.User, 0), err
	}

	for rows.Next() {
		var userID int64
		rows.Scan(&userID)

		singleUser, _, err := DumbGetUser(userID)
		if err != nil {
			return make([]models.User, 0), err
		}

		users = append(users, singleUser)
	}

	return users, nil
}

//PatchUser patches the given user. It DOES NOT update permissions and the password
func PatchUser(User models.User) error {
	_, err := Connection.Exec(`UPDATE "Users" SET "Username" = $1, "Firstname" = $2, "Lastname" = $3, "Mail" = $4 WHERE "ID" = $5`, User.Username, User.Firstname, User.Lastname, User.Mail, User.ID)
	return err
}

//DeactivateUser deactivates the user in the database
func DeactivateUser(UserID int64) error {
	_, err := Connection.Exec(`UPDATE "Users" SET "Active" = false WHERE "ID" = $1`, UserID)
	return err
}

//UpdatePassword updates the current password of the user
func UpdatePassword(UserID int64, Hash string) error {
	_, err := Connection.Exec(`UPDATE "Users" SET "Password" = $1 WHERE "ID" = $2`, Hash, UserID)
	return err
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

//GetTicketsOfUser returns all tickets where the user is the owner
func GetTicketsOfUser(UserID int64) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	rows, err := Connection.Query(`SELECT "ID" FROM "Tickets" WHERE "Owner" = $1`, UserID)
	if err != nil {
		return tickets, err
	}

	for rows.Next() {
		var id int64
		rows.Scan(&id)
		ticket, _, err := GetTicket(id, false)
		if err != nil {
			dev.LogError(err, "Couldn't get ticket: "+err.Error())
		}

		tickets = append(tickets, ticket)
	}

	return tickets, err
}
