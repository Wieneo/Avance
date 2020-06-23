package db

import (
	"database/sql"
	"errors"
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetRecipients returns all recipients assigned to the ticket
func GetRecipients(TicketID int64) (models.RecipientCollection, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting ALL recipients for ticket %d", TicketID))
	var recipients models.RecipientCollection

	dev.LogDebug(fmt.Sprintf("[DB] Initializing empty recipient arrays"))
	recipients.Admins = make([]models.Recipient, 0)
	recipients.Requestors = make([]models.Recipient, 0)
	recipients.Readers = make([]models.Recipient, 0)

	rows, err := Connection.Query(`SELECT "ID", "Type", "User", "Mail" FROM "Recipients" WHERE "Ticket" = $1`, TicketID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] An error happened while getting recipients for ticket %d: %s", TicketID, err.Error()))
		return recipients, err
	}

	for rows.Next() {
		var singleRecipient models.Recipient
		var userID sql.NullInt64
		var mail sql.NullString
		var rType models.RecipientType
		rows.Scan(&singleRecipient.ID, &rType, &userID, &mail)

		if userID.Valid {
			dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Resolving known user recipient: %d", TicketID, userID.Int64))
			singleRecipient.User.Valid = true
			singleRecipient.User.Value, _, err = getUser(userID.Int64, false)
			if err != nil {
				dev.LogDebug(fmt.Sprintf("[DB] [T: %d] User %d lookup failed: %s", TicketID, userID.Int64, err.Error()))
				return models.RecipientCollection{}, err
			}
		} else if mail.Valid {
			singleRecipient.Mail = mail.String
		} else {
			dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Got invalid recipient: %d", TicketID, singleRecipient.ID))
			return models.RecipientCollection{}, errors.New("Recipient is invalid! (User and Mail empty)")
		}

		switch rType {
		case models.Requestors:
			{
				recipients.Requestors = append(recipients.Requestors, singleRecipient)
				break
			}
		case models.Readers:
			{
				recipients.Readers = append(recipients.Readers, singleRecipient)
				break
			}
		case models.Admins:
			{
				recipients.Admins = append(recipients.Admins, singleRecipient)
				break
			}
		}
	}

	dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Got the following recipients: Admins: %d, Requestors: %d, Readers: %d", TicketID, len(recipients.Admins), len(recipients.Requestors), len(recipients.Readers)))

	return recipients, nil
}

//AddUserRecipient appends a new recipient in form of a existing user to the ticket
func AddUserRecipient(TicketID, UserID int64, Type models.RecipientType) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Adding known user recipient: %d (Type: %s)", TicketID, UserID, Type.String()))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Recipients" ("Type", "User", "Ticket") VALUES ($1, $2, $3) RETURNING "ID"`, Type, UserID, TicketID).Scan(&newID)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Recipient %d was added", TicketID, newID))
		err = TicketWasModified(TicketID)
	}

	return newID, err
}

//AddMailRecipient appends a new recipient in form of a mail address
func AddMailRecipient(TicketID int64, Mail string, Type models.RecipientType) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Adding mail recipient: %s (Type: %s)", TicketID, Mail, Type.String()))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Recipients" ("Type", "Mail", "Ticket") VALUES ($1, $2, $3) RETURNING "ID"`, Type, Mail, TicketID).Scan(&newID)

	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] [T: %d] Recipient %d was added", TicketID, newID))
		err = TicketWasModified(TicketID)
	}

	return newID, err
}

//RemoveRecipient removes a specified recipient
func RemoveRecipient(RecipientID int64) error {
	dev.LogDebug(fmt.Sprintf("[DB] Removing recipient %d", RecipientID))
	_, err := Connection.Exec(`DELETE FROM "Recipients" WHERE "ID" = $1`, RecipientID)
	dev.LogDebug(fmt.Sprintf("[DB] Removed recipient %d", RecipientID))
	return err
}
