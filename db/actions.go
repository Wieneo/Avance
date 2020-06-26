package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetActions returns all actions associated with a ticket
func GetActions(TicketID int64) ([]models.Action, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting Actions for ticket %d", TicketID))
	actions := make([]models.Action, 0)
	rows, err := Connection.Query(`SELECT "ID", "Type", "Title", "Content", "IssuedAt", "IssuedBy", "Tasks" FROM "Actions" WHERE "Ticket" = $1 ORDER BY "ID" DESC`, TicketID)
	if err != nil {
		return actions, err
	}

	cachedUsers := make(map[int64]models.User, 0)

	for rows.Next() {
		var singleAction models.Action
		var rawUserID sql.NullInt64
		var rawTasks sql.NullString
		rows.Scan(&singleAction.ID, &singleAction.Type, &singleAction.Title, &singleAction.Content, &singleAction.IssuedAt, &rawUserID, &rawTasks)

		if rawUserID.Valid {

			//If user is not cached
			if _, found := cachedUsers[rawUserID.Int64]; !found {
				user, _, err := GetUser(rawUserID.Int64)
				if err != nil {
					return make([]models.Action, 0), err
				}

				cachedUsers[rawUserID.Int64] = user
			}

			singleAction.IssuedBy.Valid = true
			singleAction.IssuedBy.Issuer = cachedUsers[rawUserID.Int64]
		}

		if rawTasks.Valid {
			err = json.Unmarshal([]byte(rawTasks.String), &singleAction.Tasks)
			if err != nil {
				return make([]models.Action, 0), err
			}
		} else {
			singleAction.Tasks = make([]int64, 0)
		}

		actions = append(actions, singleAction)
	}
	dev.LogDebug(fmt.Sprintf("[DB] Got %d actions for ticket %d", len(actions), TicketID))

	return actions, nil
}

//AddAction adds an action to a ticket
func AddAction(TicketID int64, Type models.ActionType, Title, Content string, IssuedBy models.Issuer) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Adding the following action to ticket %d: Type: %d, Title: %s, Content: Omitted", TicketID, Type, Title))
	var newID int64

	var issuedByReal sql.NullInt64
	if IssuedBy.Valid {
		issuedByReal.Valid = true
		issuedByReal.Int64 = IssuedBy.Issuer.ID
	}

	err := Connection.QueryRow(`INSERT INTO "Actions" ("Type", "Title", "Content", "Ticket", "IssuedAt", "IssuedBy") VALUES ($1,$2,$3,$4,$5,$6) RETURNING "ID"`, Type, Title, Content, TicketID, time.Now(), issuedByReal).Scan(&newID)

	if err == nil {
		err = TicketWasModified(TicketID)
		dev.LogDebug(fmt.Sprintf("[DB] Created action %d for ticket %d", newID, TicketID))
	}

	ticket, _, err := GetTicketUnsafe(TicketID, models.WantedProperties{Recipients: true})
	if err == nil {
		dev.LogDebug(fmt.Sprintf("Preparing Notifications for ticket %d", TicketID))
		go QueueActionNotification(ticket, models.Action{
			ID:       newID,
			Title:    Title,
			Content:  Content,
			Type:     Type,
			IssuedBy: IssuedBy,
			IssuedAt: time.Now(),
		})
	}

	return newID, err
}
