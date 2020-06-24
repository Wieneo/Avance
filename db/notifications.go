package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//QueueActionNotification send all needed notifications into the queue
func QueueActionNotification(Ticket models.Ticket, Action models.Action) error {
	dev.LogDebug(fmt.Sprintf("[DB] Queueing notification for ticket %d", Ticket.ID))
	var Error error
	dev.LogDebug(fmt.Sprintf("[DB] Retrieving all recipients for ticket %d", Ticket.ID))
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
	dev.LogDebug(fmt.Sprintf("[DB] Queueing Mail Notification for ticket %d", Ticket.ID))

	var trueRecipient string
	if Recipient.User.Valid {
		dev.LogDebug(fmt.Sprintf("[DB] Recipient is user %d", Recipient.User.Value.ID))
		trueRecipient = strconv.FormatInt(Recipient.User.Value.ID, 10)
	} else {
		dev.LogDebug(fmt.Sprintf("[DB] Recipient is mail %s", Recipient.Mail))
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

		dev.LogDebug(fmt.Sprintf("[DB] Parsing previous notifications"))

		err := json.Unmarshal([]byte(rawOldNotifications), &oldNotifications)
		if err != nil {
			dev.LogError(err, fmt.Sprintf("Notification Task for %d seems to be corrupt!", Recipient.ID))
			//ToDo: Better Error Handling -> Reject Request / Remove corrupt notification task
			return err
		}

		oldNotifications.Notifications = append(oldNotifications.Notifications, models.Notification{
			Title:   Action.Title,
			Content: Action.Content,
			Action: struct {
				Valid bool
				Value models.Action
			}{
				Valid: true,
				Value: Action,
			},
		})

		dev.LogDebug(fmt.Sprintf("[DB] Mahrshalling new notificaitons struct for ticket %d", Ticket.ID))

		rawJSON, _ := json.Marshal(oldNotifications)

		dev.LogDebug(fmt.Sprintf("[DB] Updating notification struct for ticket %d", Ticket.ID))

		//Append new notification in database
		if _, err := Connection.Exec(`UPDATE "Tasks" SET "Task" = $1 WHERE "ID" = $2`, string(rawJSON), oldTaskID); err != nil {
			return err
		}

		dev.LogDebug(fmt.Sprintf("Task %d expanded to %d notifications", oldTaskID, len(oldNotifications.Notifications)))

		err = addTasksToAction(oldTaskID, Action)
		if err != nil {
			return err
		}

	} else {
		//We need to create a brand new notification task
		rows.Close()
		dev.LogDebug(fmt.Sprintf(`Found no preceeding notification for recipient %d -> Creating new task`, Recipient.ID))
		var notifications models.NotificationCollection
		notifications.Subject = fmt.Sprintf("Update about ticket %d", Ticket.ID) //ToDo: Make more dynamic later? Maybe we want to set other subjects then "UPDATE" -> Ticket creation?
		//Later used by worker to determine what to do
		notifications.NotifyType = models.Mail
		notifications.Notifications = append(notifications.Notifications, models.Notification{
			Title:   Action.Title,
			Content: Action.Content,
			Action: struct {
				Valid bool
				Value models.Action
			}{
				Valid: true,
				Value: Action,
			},
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

		dev.LogDebug(fmt.Sprintf("[DB] Struct is ready to be dumped to the database for ticket %d", Ticket.ID))
		dev.LogDebug(fmt.Sprintf("[DB] Creating task for notification for ticket %d", Ticket.ID))

		if newTaskID, err := CreateTask(models.SendNotification, string(rawJSON), Interval, sql.NullString{Valid: true, String: trueRecipient}, sql.NullInt64{Valid: true, Int64: Ticket.ID}, sql.NullInt32{Valid: true, Int32: int32(models.Mail)}); err == nil {
			err = addTasksToAction(newTaskID, Action)
			if err != nil {
				return err
			}
		} else {
			return err
		}

	}

	return nil
}

func addTasksToAction(TaskID int64, Action models.Action) error {
	//As the Action struct is only temporary we get the current tasks from the database
	var rawJSON sql.NullString
	dev.LogDebug(fmt.Sprintf("[DB] Retrieving previous tasks for action %d", Action.ID))
	err := Connection.QueryRow(`SELECT "Tasks" FROM "Actions" WHERE "ID" = $1`, Action.ID).Scan(&rawJSON)

	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while retrieving previous tasks for action %d: %s", Action.ID, err.Error()))
		return err
	}

	tasks := make([]int64, 0)
	if rawJSON.Valid {
		err = json.Unmarshal([]byte(rawJSON.String), &tasks)
	}

	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Parsing previous tasks for action %d failed: %s", Action.ID, err.Error()))
		return err
	}

	tasks = append(tasks, TaskID)

	newJSON, _ := json.Marshal(tasks)

	dev.LogDebug(fmt.Sprintf("[DB] Expanding tasks of action %d to %d tasks", Action.ID, len(tasks)))
	Connection.Exec(`UPDATE "Actions" SET "Tasks" = $1 WHERE "ID" = $2`, string(newJSON), Action.ID)
	return nil
}
