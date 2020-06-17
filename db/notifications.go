package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//QueueActionNotification send all needed notifications into the queue
func QueueActionNotification(Ticket models.Ticket, Action models.Action) error {
	var Error error
	for _, k := range Ticket.AllRecipients() {
		//Admins get all
		if k.Type == models.Admins ||
			//Answer to requestors and readers
			((k.Type == models.Requestors || k.Type == models.Readers) && Action.Type == models.Answer) {
			dev.LogDebug(fmt.Sprintf("Sending notification to %d (Local-ID)", k.ID))

			/*Currently we assume that there is only E-Mail as notification type
			This should probably later changed to something more dynamic (numbers / handles for telegram, etc...)
			*/

			//Check if user / unknown e-mail was specified
			if k.User.Valid {
				Error = GetSettings(&k.User.Value)
				if Error != nil {
					dev.LogError(Error, fmt.Sprintf("Couldn't get notification settings for task: %s", Error.Error()))
					return Error
				}
				if k.User.Value.Settings.Notification.MailNotificationEnabled && k.User.Value.Settings.Notification.MailNotificationAboutUpdates {
					Error = sendMailActionNotificationIntoQueue(Ticket, Action, k)
				}
			} else {
				Error = sendMailActionNotificationIntoQueue(Ticket, Action, k)
			}

			if Error != nil {
				dev.LogError(Error, fmt.Sprintf("Couldn't queue notification task: %s", Error.Error()))
			}
		}
	}
	return Error
}

func sendMailActionNotificationIntoQueue(Ticket models.Ticket, Action models.Action, Recipient models.Recipient) error {
	var trueRecipient string
	if Recipient.User.Valid {
		trueRecipient = strconv.FormatInt(Recipient.User.Value.ID, 10)
	} else {
		trueRecipient = Recipient.Mail
	}

	//Check if there already is a notification task with the e-mail type
	rows, err := Connection.Query(`SELECT "ID", "Task" FROM "Tasks" WHERE "Ticket" = $1 AND "Recipient" = $2 AND "Status" = 0 AND "NotifiationType" = $3`, Ticket.ID, trueRecipient, models.Mail)
	if err != nil {
		return err
	}

	//We expand the existing task
	if rows.Next() {
		//GetOldTask
		var oldTaskID int64
		var oldNotifications models.NotificationCollection
		var rawOldNotifications string
		rows.Scan(&oldTaskID, &rawOldNotifications)
		rows.Close()

		err := json.Unmarshal([]byte(rawOldNotifications), &oldNotifications)
		if err != nil {
			dev.LogError(err, fmt.Sprintf("Notification Task for %d seems to be corrupt!", Recipient.ID))
			//ToDo: Better Error Handling -> Reject Request / Remove corrupt notification task
			return err
		}

		oldNotifications.Notifications = append(oldNotifications.Notifications, models.Notification{
			Title:   Action.Title,
			Content: Action.Content,
		})

		rawJSON, _ := json.Marshal(oldNotifications)

		//Append new notification in database
		if _, err := Connection.Exec(`UPDATE "Tasks" SET "Task" = $1 WHERE "ID" = $2`, string(rawJSON), oldTaskID); err != nil {
			return err
		}

		dev.LogDebug(fmt.Sprintf("Task %d expanded to %d notifications", oldTaskID, len(oldNotifications.Notifications)))

	} else {
		//We need to create a brand new notification task
		rows.Close()
		dev.LogDebug(fmt.Sprintf(`Found no preceeding notification for recipient %d -> Creating new task`, Recipient.ID))
		var notifications models.NotificationCollection
		notifications.Subject = fmt.Sprintf("Update about ticket %d", Ticket.ID) //ToDo: Make more dynamic later? Maybe we want to set other subjects then "UPDATE" -> Ticket creation?
		//Later used be worker to determine what to do
		notifications.NotifyType = models.Mail
		notifications.Notifications = append(notifications.Notifications, models.Notification{
			Title:   Action.Title,
			Content: Action.Content,
		})
		rawJSON, _ := json.Marshal(notifications)

		//Interval is always true and is set to the users preference / the default (30 Sec.) -> ToDo: Maybe make that configurable via the general settings?
		Interval := sql.NullInt32{
			Valid: true,
		}

		if Recipient.User.Valid {
			err := GetSettings(&Recipient.User.Value)
			if err != nil {
				dev.LogError(err, fmt.Sprintf(`Couldn't get settings of user %d: %s`, Recipient.User.Value.ID, err.Error()))
				return err
			}

			Interval.Int32 = int32(Recipient.User.Value.Settings.Notification.MailNotificationFrequency)
		} else {
			Interval.Int32 = 30
		}

		if _, err := CreateTask(models.SendNotification, string(rawJSON), Interval, sql.NullString{Valid: true, String: trueRecipient}, sql.NullInt64{Valid: true, Int64: Ticket.ID}, sql.NullInt32{Valid: true, Int32: int32(models.Mail)}); err != nil {
			return err
		}
	}

	return nil
}
