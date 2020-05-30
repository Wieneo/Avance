package db

import (
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetActions returns all actions associated with a ticket
func GetActions(TicketID int64) ([]models.Action, error) {
	actions := make([]models.Action, 0)
	rows, err := Connection.Query(`SELECT "ID", "Type", "Title", "Content", "IssuedAt", "IssuedBy" FROM "Actions" WHERE "Ticket" = $1 ORDER BY "ID" DESC`, TicketID)
	if err != nil {
		return actions, err
	}

	for rows.Next() {
		var singleAction models.Action
		var rawUserID int64
		rows.Scan(&singleAction.ID, &singleAction.Type, &singleAction.Title, &singleAction.Content, &singleAction.IssuedAt, &rawUserID)

		user, _, err := GetUser(rawUserID)
		if err != nil {
			return make([]models.Action, 0), err
		}

		singleAction.IssuedBy = user
		actions = append(actions, singleAction)
	}

	return actions, nil
}

//AddAction adds an action to a ticket
func AddAction(TicketID int64, Type models.ActionType, Title, Content string, IssuedBy int64) (int64, error) {
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Actions" ("Type", "Title", "Content", "Ticket", "IssuedAt", "IssuedBy") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "ID"`, Type, Title, Content, TicketID, time.Now(), IssuedBy).Scan(&newID)
	return newID, err
}
