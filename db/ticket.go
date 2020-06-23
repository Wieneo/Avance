package db

import (
	"database/sql"
	"fmt"
	"time"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetTicketUnsafe ignores if the queue matches the tickets queue
func GetTicketUnsafe(TicketID int64, ResolveRelations bool) (models.Ticket, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting queue for ticket %d in UNSAFE mode (Resolving relations: %t)", TicketID, ResolveRelations))

	var queueID int64
	err := Connection.QueryRow(`SELECT "Queue" FROM "Tickets" WHERE "ID" = $1`, TicketID).Scan(&queueID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			dev.LogDebug(fmt.Sprintf("[DB] Ticket %d wasn't found", TicketID))
			return models.Ticket{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Ticket{}, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Reverse lookup for queue %d succeded", TicketID, queueID))
	return GetTicket(TicketID, queueID, ResolveRelations)
}

//GetTicket returns the populated ticket struct from models
func GetTicket(TicketID int64, QueueID int64, ResolveRelations bool) (models.Ticket, bool, error) {
	started := time.Now()
	dev.LogDebug(fmt.Sprintf("[DB] Getting ticket %d in queue %d (Resolving relations: %t)", TicketID, QueueID, ResolveRelations))

	var ticket models.Ticket
	err := Connection.QueryRow(`SELECT "t"."ID", "t"."Title", "t"."Description", "t"."Queue" AS "QueueID", "t"."Owner" AS "OwnerID", "t"."Severity" AS "SeverityID", "t"."Status" AS "StatusID","t"."CreatedAt","t"."LastModified","t"."StalledUntil","t"."Meta" FROM "Tickets" AS "t" WHERE "ID" = $1 AND "Queue" = $2`, TicketID, QueueID).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.QueueID, &ticket.OwnerID, &ticket.SeverityID, &ticket.StatusID, &ticket.CreatedAt, &ticket.LastModified, &ticket.StalledUntil, &ticket.Meta)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			dev.LogDebug(fmt.Sprintf("[DB] Ticket %d wasn't found in queue %d", TicketID, QueueID))
			return models.Ticket{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Ticket{}, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Looking up queue for ticket %d", TicketID))
	ticket.Queue, _, err = GetQueueUNSAFE(ticket.QueueID)
	if ticket.OwnerID.Valid {
		dev.LogDebug(fmt.Sprintf("[DB] Looking up owner for ticket %d", TicketID))
		ticket.Owner, _, err = DumbGetUser(ticket.OwnerID.Int64)
	}
	dev.LogDebug(fmt.Sprintf("[DB] Looking up severity for ticket %d", TicketID))
	ticket.Severity, _, err = GetSeverityUNSAFE(ticket.SeverityID)
	dev.LogDebug(fmt.Sprintf("[DB] Looking up status for ticket %d", TicketID))
	ticket.Status, _, err = GetStatusUNSAFE(ticket.StatusID)

	if ResolveRelations {
		dev.LogDebug(fmt.Sprintf("[DB] Looking up relations for ticket %d", TicketID))
		ticket.Relations, err = GetTicketRelations(TicketID)
		if err != nil {
			dev.LogDebug(fmt.Sprintf("[DB] Error while looking up relations for ticket %d: %s", TicketID, err.Error()))
			return models.Ticket{}, true, err
		}
	} else {
		//Prevent Relations from being NULL
		ticket.Relations = make([]models.Relation, 0)
	}

	dev.LogDebug(fmt.Sprintf("[DB] Looking up actions for ticket %d", TicketID))
	ticket.Actions, err = GetActions(ticket.ID)
	if err != nil {
		return models.Ticket{}, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Looking up recipients for ticket %d", TicketID))
	ticket.Recipients, err = GetRecipients(ticket.ID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error while looking up recipients for ticket %d: %s", TicketID, err.Error()))
		return models.Ticket{}, true, err
	}

	end := time.Now()

	dev.LogDebug(fmt.Sprintf("[DB] Retrieving ticket %d took %s", TicketID, end.Sub(started)))
	return ticket, true, nil
}

//CreateTicket creates a ticket and returns the new id
func CreateTicket(Title string, Description string, Queue int64, OwnedByNobody bool, Owner int64, Severity models.Severity, Status models.Status, IsStalled bool, StalledUntil string) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Creating ticket in queue %d with values: Title: %s, Owner: %d, Severity: %d, Status: %d", Queue, Title, Owner, Severity.ID, Status.ID))
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
		dev.LogDebug(fmt.Sprintf("[DB] Created ticket %d -> Adding initial action", newID))
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

		if err == nil {
			dev.LogDebug(fmt.Sprintf("[DB] Initial Action added to ticket %d", newID))
		}
	}
	return newID, err
}

//PatchTicket patches the given ticket
func PatchTicket(Ticket models.Ticket) (models.Ticket, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Patching ticket %d", Ticket.ID))
	_, err := Connection.Exec(`UPDATE "Tickets" SET "Title" = $1, "Description" = $2, "Queue" = $3, "Owner" = $4, "Severity" = $5, "Status" = $6, "StalledUntil" = $7, "Meta" = $8 WHERE "ID" = $9`, Ticket.Title, Ticket.Description, Ticket.QueueID, Ticket.OwnerID, Ticket.SeverityID, Ticket.StatusID, Ticket.StalledUntil, Ticket.Meta, Ticket.ID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while patching ticket %d: %s", Ticket.ID, err.Error()))
		return models.Ticket{}, err
	}
	if TicketWasModified(Ticket.ID) != nil {
		return models.Ticket{}, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Retrieving new ticket %d after patching", Ticket.ID))
	ticket, _, err := GetTicketUnsafe(Ticket.ID, true)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] [1] Error happened while patching ticket %d: %s", Ticket.ID, err.Error()))
		return models.Ticket{}, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Patched ticket %d", Ticket.ID))
	return ticket, err
}

//TicketWasModified refreshes the last modified column of a ticket
func TicketWasModified(TicketID int64) error {
	dev.LogDebug(fmt.Sprintf("[DB] Updating LastModified value for ticket %d", TicketID))
	_, err := Connection.Exec(`UPDATE "Tickets" SET "LastModified" = $1 WHERE "ID" = $2`, time.Now(), TicketID)
	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Succesfully updated LastModified value for ticket %d", TicketID))
	}
	return err
}

//GetTicketsInQueue returns all tickets in a give queue
func GetTicketsInQueue(QueueID int64, ShowInvisible bool) ([]models.Ticket, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting ALL tickets in queue %d (Showing invisible: %t)", QueueID, ShowInvisible))

	tickets := make([]models.Ticket, 0)
	rows, err := Connection.Query(`SELECT "t"."ID", "t"."Title", "t"."Description", "t"."Queue" AS "QueueID", "t"."Owner" AS "OwnerID", "t"."Severity" AS "SeverityID", "t"."Status" AS "StatusID","t"."CreatedAt","t"."LastModified","t"."StalledUntil","t"."Meta" FROM "Tickets" AS "t" WHERE "Queue" = $1`, QueueID)
	if err != nil {
		dev.LogError(err, err.Error())
		return make([]models.Ticket, 0), err
	}

	for rows.Next() {
		var ticket models.Ticket
		rows.Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.QueueID, &ticket.OwnerID, &ticket.SeverityID, &ticket.StatusID, &ticket.CreatedAt, &ticket.LastModified, &ticket.StalledUntil, &ticket.Meta)
		dev.LogDebug(fmt.Sprintf("[DB] Retrieving ticket %d in queue %d", ticket.ID, QueueID))
		dev.LogDebug(fmt.Sprintf("[DB] Retrieving queue for ticket %d in queue %d", ticket.ID, QueueID))
		ticket.Queue, _, err = GetQueueUNSAFE(ticket.QueueID)
		if ticket.OwnerID.Valid {
			dev.LogDebug(fmt.Sprintf("[DB] Retrieving owner for ticket %d in queue %d", ticket.ID, QueueID))
			ticket.Owner, _, err = GetUser(ticket.OwnerID.Int64)
		}
		dev.LogDebug(fmt.Sprintf("[DB] Retrieving severity for ticket %d in queue %d", ticket.ID, QueueID))
		ticket.Severity, _, err = GetSeverityUNSAFE(ticket.SeverityID)
		dev.LogDebug(fmt.Sprintf("[DB] Retrieving status for ticket %d in queue %d", ticket.ID, QueueID))
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

	dev.LogDebug(fmt.Sprintf("[DB] Got %d tickets in queue %d", len(tickets), QueueID))
	return tickets, nil
}
