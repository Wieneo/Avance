package db

import (
	"database/sql"
	"errors"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetRecipients returns all recipients assigned to the ticket
func GetRecipients(TicketID int64) (models.RecipientCollection, error) {
	var recipients models.RecipientCollection

	recipients.Admins = make([]models.Recipient, 0)
	recipients.Requestors = make([]models.Recipient, 0)
	recipients.Readers = make([]models.Recipient, 0)

	rows, err := Connection.Query(`SELECT "ID", "Type", "User", "Mail" FROM "Recipients" WHERE "Ticket" = $1`, TicketID)
	if err != nil {
		return recipients, err
	}

	for rows.Next() {
		var singleRecipient models.Recipient
		var userID sql.NullInt64
		var mail sql.NullString
		var rType models.RecipientType
		rows.Scan(&singleRecipient.ID, &rType, &userID, &mail)

		if userID.Valid {
			singleRecipient.User.Valid = true
			singleRecipient.User.Value, _, err = getUser(userID.Int64, false)
			if err != nil {
				return models.RecipientCollection{}, err
			}
		} else if mail.Valid {
			singleRecipient.Mail = mail.String
		} else {
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

	return recipients, nil
}
