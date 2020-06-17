package functions

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//SendNotifications send the notifications specified in the notification task
func SendNotifications(Task models.WorkerTask) error {
	var Notifications models.NotificationCollection
	err := json.Unmarshal([]byte(Task.Data), &Notifications)
	if err != nil {
		return err
	}

	switch Notifications.NotifyType {
	case models.Mail:
		return sendMailNotification(Task, Notifications)
	}

	return errors.New("No action was taken")
}

func sendMailNotification(Task models.WorkerTask, Notifications models.NotificationCollection) error {
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
			return err
		}

		if !found {
			dev.LogWarn(fmt.Sprintf(`User %d wasn't found anymore -> Skipping notification`, userid))
			return nil
		}

		realRecipient = user.Mail
	}

	//Reference https://gist.github.com/chrisgillis/10888032

	var content string

	//ToDo: Templating System for E-Mail notifications
	for _, k := range Notifications.Notifications {
		content += "<h3>" + k.Title + "</h3>"
		content += k.Content
		content += "<br>"
	}

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = config.CurrentConfig.SMTP.From
	headers["To"] = realRecipient
	headers["Subject"] = "[Avance] " + Notifications.Subject

	// Setup message
	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	//Mime needs to be appended last in order to make the headers before it work
	message += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\r\n" + content

	auth := smtp.PlainAuth("", config.CurrentConfig.SMTP.User, config.CurrentConfig.SMTP.Password, config.CurrentConfig.SMTP.Host)
	// TLS config
	tlsconfig := &tls.Config{
		ServerName: config.CurrentConfig.SMTP.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", config.CurrentConfig.SMTP.Host, config.CurrentConfig.SMTP.Port), tlsconfig)
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	c, err := smtp.NewClient(conn, config.CurrentConfig.SMTP.Host)
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	c.StartTLS(tlsconfig)

	// Auth
	if err = c.Auth(auth); err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	// To && From
	if err = c.Mail(config.CurrentConfig.SMTP.From); err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	if err = c.Rcpt(realRecipient); err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	err = w.Close()
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification to %s: %s", realRecipient, err.Error()))
		return err
	}

	c.Quit()

	return nil
}
