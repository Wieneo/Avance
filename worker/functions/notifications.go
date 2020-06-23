package functions

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gitlab.gnaucke.dev/avance/avance-app/v2/config"
	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/worker/smtp"
)

//SendNotifications send the notifications specified in the notification task
func SendNotifications(Task models.WorkerTask) (bool, error) {
	var Notifications models.NotificationCollection
	err := json.Unmarshal([]byte(Task.Data), &Notifications)
	if err != nil {
		return false, err
	}

	switch Notifications.NotifyType {
	case models.Mail:
		return sendMailNotification(Task, Notifications)
	}

	return false, errors.New("No action was taken")
}

func sendMailNotification(Task models.WorkerTask, Notifications models.NotificationCollection) (bool, error) {
	if !config.CurrentConfig.SMTP.Enabled {
		err := errors.New("Mail Notification without SMTP")
		dev.LogError(err, "Mail Notification was queued without having SMTP configured! Ignoring Task!")
		return false, err
	}

	var realRecipient string
	//Task.Recipient is either an ID (the user-id) OR a string containing the unknown e-mail address
	//Trying to check which one of the options was specififed
	if userid, err := strconv.ParseInt(Task.Recipient.String, 10, 64); err != nil {
		//Unknown E-mail
		realRecipient = Task.Recipient.String
	} else {
		user, found, err := db.GetUser(userid)
		if err != nil {
			dev.LogError(err, fmt.Sprintf("Can't get user %d: %s", userid, err.Error()))
			return false, err
		}

		if !found {
			dev.LogWarn(fmt.Sprintf(`User %d wasn't found anymore -> Skipping notification`, userid))
			return false, nil
		}

		realRecipient = user.Mail
	}

	var content string

	//ToDo: Templating System for E-Mail notifications -> Load General E-Mail Template
	for _, k := range Notifications.Notifications {
		//Check if we send ticket notifications
		if k.Action.Valid {
			//ToDo: Load Single Templates for Action - Notifications
			content += "<h3>" + k.Title + "</h3>"
			content += k.Content
			content += "<br>"
		} else {
			content += "<h3>" + k.Title + "</h3>"
			content += k.Content
			content += "<br>"
		}
	}

	err := smtp.SendMail(realRecipient, Notifications.Subject, content)
	if err != nil {
		return true, err
	}

	return false, nil
}
