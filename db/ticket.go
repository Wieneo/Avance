package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetTicket returns the populated ticket struct from models
func GetTicket(TicketID int64) (models.Ticket, bool, error) {
	var ticket models.Ticket
	err := Connection.QueryRow(`SELECT "t"."ID", "t"."Title", "t"."Description", "t"."Queue" AS "QueueID", "t"."Owner" AS "OwnerID", "t"."Severity" AS "SeverityID", "t"."Status" AS "StatusID","t"."CreatedAt","t"."LastModified","t"."StalledUntil","t"."Meta" FROM "Tickets" AS "t" WHERE "ID" = $1`, TicketID).Scan(&ticket.ID, &ticket.Title, &ticket.Description, &ticket.QueueID, &ticket.OwnerID, &ticket.SeverityID, &ticket.StatusID, &ticket.CreatedAt, &ticket.LastModified, &ticket.StalledUntil, &ticket.Meta)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Ticket{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Ticket{}, true, err
	}

	ticket.Queue, _, err = GetQueueUNSAFE(ticket.QueueID)
	if ticket.OwnerID.Valid {
		ticket.Owner, err = GetUser(ticket.OwnerID.Int64)
	}
	ticket.Severity, _, err = GetSeverityUNSAFE(ticket.SeverityID)
	ticket.Status, _, err = GetStatusUNSAFE(ticket.StatusID)
	return ticket, true, nil
}

//GetTicketsInQueue returns all tickets in a give queue
func GetTicketsInQueue(QueueID int64) ([]models.Ticket, error) {
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
			ticket.Owner, err = GetUser(ticket.OwnerID.Int64)
		}
		ticket.Severity, _, err = GetSeverityUNSAFE(ticket.SeverityID)
		ticket.Status, _, err = GetStatusUNSAFE(ticket.StatusID)

		if err != nil {
			dev.LogError(err, err.Error())
			return make([]models.Ticket, 0), err
		}
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}
