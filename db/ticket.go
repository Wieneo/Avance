package db

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetTicketUnsafe ignores if the queue matches the tickets queue
func GetTicketUnsafe(TicketID int64, ResolveRelations bool) (models.Ticket, bool, error) {
	var queueID int64
	err := Connection.QueryRow(`SELECT "Queue" FROM "Tickets" WHERE "ID" = $1`, TicketID).Scan(&queueID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Ticket{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Ticket{}, true, err
	}

	return GetTicket(TicketID, queueID, ResolveRelations)
}

//GetTicket returns the populated ticket struct from models
func GetTicket(TicketID int64, QueueID int64, ResolveRelations bool) (models.Ticket, bool, error) {
	var ticket models.Ticket
	err := Connection.QueryRow(`SELECT "t"."ID", "t"."Title", "t"."Description", "t"."Queue" AS "QueueID", "t"."Owner" AS "OwnerID", "t"."Severity" AS "SeverityID", "t"."Status" AS "StatusID","t"."CreatedAt","t"."LastModified","t"."StalledUntil","t"."Meta" FROM "Tickets" AS "t" WHERE "ID" = $1 AND "Queue" = $2`, TicketID, QueueID).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.QueueID, &ticket.OwnerID, &ticket.SeverityID, &ticket.StatusID, &ticket.CreatedAt, &ticket.LastModified, &ticket.StalledUntil, &ticket.Meta)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Ticket{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Ticket{}, true, err
	}

	ticket.Queue, _, err = GetQueueUNSAFE(ticket.QueueID)
	if ticket.OwnerID.Valid {
		ticket.Owner, _, err = DumbGetUser(ticket.OwnerID.Int64)
	}
	ticket.Severity, _, err = GetSeverityUNSAFE(ticket.SeverityID)
	ticket.Status, _, err = GetStatusUNSAFE(ticket.StatusID)

	if ResolveRelations {
		ticket.Relations, err = GetTicketRelations(TicketID)
		if err != nil {
			return models.Ticket{}, true, err
		}
	} else {
		//Prevent Relations from being NULL
		ticket.Relations = make([]models.Relation, 0)
	}

	ticket.Actions, err = GetActions(ticket.ID)
	if err != nil {
		return models.Ticket{}, true, err
	}

	ticket.Recipients, err = GetRecipients(ticket.ID)
	if err != nil {
		return models.Ticket{}, true, err
	}

	return ticket, true, nil
}

//CreateTicket creates a ticket and returns the new id
func CreateTicket(Title string, Description string, Queue int64, OwnedByNobody bool, Owner int64, Severity models.Severity, Status models.Status, IsStalled bool, StalledUntil string) (int64, error) {
	var newID int64
	var trueOwner sql.NullInt64
	var trueStall sql.NullString

	if !OwnedByNobody {
		trueOwner.Valid = true
		trueOwner.Int64 = Owner
	}

	if IsStalled {
		trueStall.Valid = true
		trueStall.String = StalledUntil
	}

	err := Connection.QueryRow(`INSERT INTO "Tickets" ("Title", "Description", "Queue", "Owner", "Severity", "Status", "CreatedAt", "LastModified", "StalledUntil", "Meta") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING "ID"`, Title, Description, Queue, trueOwner, Severity.ID, Status.ID, time.Now(), time.Now(), trueStall, "{}").Scan(&newID)

	if err == nil {
		resolvedQueue, _, _ := GetQueueUNSAFE(Queue)
		resolvedOwner := ""
		if !OwnedByNobody {
			owner, _, _ := GetUser(Owner)
			resolvedOwner = owner.Username
		}

		//Add "Ticket Created" Action from user System
		_, err = AddAction(newID, models.Unspecific, "Ticket Created", fmt.Sprintf(`Ticket was created with the following properties:<br>
			<ul>
				<li>Title: <i>%s</i></li>
				<li>Description: <i>%s</i></li>
				<li>Queue: <i>%s</i></li>
				<li>Owner: <i>%s</i></li>
				<li>Severity: <i>%s</i></li>
				<li>Status: <i>%s</i></li>
				<li>StalledUntil: <i>%s</i></li>
			</ul>`, Title, Description, resolvedQueue.Name, resolvedOwner, Severity.Name, Status.Name, StalledUntil), models.Issuer{Valid: false})
	}
	return newID, err
}

//PatchTicket patches the given ticket
func PatchTicket(Ticket models.Ticket) (models.Ticket, error) {
	_, err := Connection.Exec(`UPDATE "Tickets" SET "Title" = $1, "Description" = $2, "Queue" = $3, "Owner" = $4, "Severity" = $5, "Status" = $6, "StalledUntil" = $7, "Meta" = $8 WHERE "ID" = $9`, Ticket.Title, Ticket.Description, Ticket.QueueID, Ticket.OwnerID, Ticket.SeverityID, Ticket.StatusID, Ticket.StalledUntil, Ticket.Meta, Ticket.ID)
	if err != nil {
		return models.Ticket{}, err
	}
	if TicketWasModified(Ticket.ID) != nil {
		return models.Ticket{}, err
	}

	ticket, _, err := GetTicketUnsafe(Ticket.ID, true)
	if err != nil {
		return models.Ticket{}, err
	}
	return ticket, err
}

//TicketWasModified refreshes the last modified column of a ticket
func TicketWasModified(TicketID int64) error {
	_, err := Connection.Exec(`UPDATE "Tickets" SET "LastModified" = $1 WHERE "ID" = $2`, time.Now(), TicketID)
	return err
}

//GetTicketsInQueue returns all tickets in a give queue
func GetTicketsInQueue(QueueID int64, ShowInvisible bool) ([]models.Ticket, error) {
	tickets := make([]models.Ticket, 0)
	rows, err := Connection.Query(`SELECT "t"."ID", "t"."Title", "t"."Description", "t"."Queue" AS "QueueID", "t"."Owner" AS "OwnerID", "t"."Severity" AS "SeverityID", "t"."Status" AS "StatusID","t"."CreatedAt","t"."LastModified","t"."StalledUntil","t"."Meta" FROM "Tickets" AS "t" WHERE "Queue" = $1`, QueueID)
	if err != nil {
		dev.LogError(err, err.Error())
		return make([]models.Ticket, 0), err
	}

	for rows.Next() {
		var ticket models.Ticket
		rows.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.QueueID, &ticket.OwnerID, &ticket.SeverityID, &ticket.StatusID, &ticket.CreatedAt, &ticket.LastModified, &ticket.StalledUntil, &ticket.Meta)
		ticket.Queue, _, err = GetQueueUNSAFE(ticket.QueueID)
		if ticket.OwnerID.Valid {
			ticket.Owner, _, err = GetUser(ticket.OwnerID.Int64)
		}
		ticket.Severity, _, err = GetSeverityUNSAFE(ticket.SeverityID)
		ticket.Status, _, err = GetStatusUNSAFE(ticket.StatusID)

		if err != nil {
			dev.LogError(err, err.Error())
			return make([]models.Ticket, 0), err
		}

		//Relations are not resolved on purpose to save on time & ressources
		//ALso i don't see the point in having relations here
		//Prevent Relations from being NULL
		ticket.Relations = make([]models.Relation, 0)

		//Actions are not resolved on to save on time & ressources
		//Prevent Relations from being NUL
		ticket.Actions = make([]models.Action, 0)

		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
