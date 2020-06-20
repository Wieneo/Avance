package smtp

import (
	"crypto/tls"
	"fmt"
	"net/smtp"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//Login connects to the SMTP server and returns the SMTP Client struct
func Login() (*smtp.Client, error) {
	auth := smtp.PlainAuth("", config.CurrentConfig.SMTP.User, config.CurrentConfig.SMTP.Password, config.CurrentConfig.SMTP.Host)
	// TLS config
	tlsconfig := &tls.Config{
		ServerName: config.CurrentConfig.SMTP.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", config.CurrentConfig.SMTP.Host, config.CurrentConfig.SMTP.Port), tlsconfig)
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't connect to SMTP Server: %s", err.Error()))
		return nil, err
	}

	c, err := smtp.NewClient(conn, config.CurrentConfig.SMTP.Host)
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't create SMTP Client Connection: %s", err.Error()))
		return nil, err
	}

	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't initiate TLS Connection to SMTP Server: %s", err.Error()))
		return nil, err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		dev.LogError(err, fmt.Sprintf("Can't authenticate to SMTP Server: %s", err.Error()))
		return nil, err
	}

	return c, nil
}

//SendMail sends a normal E-Mail to the specified recipient. Login is done automatically and must not be called beforehand.
func SendMail(Recipient, Subject, Message string) error {
	c, err := Login()
	if err != nil {
		return err
	}

	//Reference https://gist.github.com/chrisgillis/10888032

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = config.CurrentConfig.SMTP.From
	headers["To"] = Recipient
	headers["Subject"] = "[Avance] " + Subject

	// Setup message
	completeMessage := ""
	for k, v := range headers {
		completeMessage += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	//Mime needs to be appended last in order to make the headers before it work
	completeMessage += "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n\r\n" + Message

	// To && From
	if err = c.Mail(config.CurrentConfig.SMTP.From); err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification [0] to %s: %s", Recipient, err.Error()))
		return err
	}

	if err = c.Rcpt(Recipient); err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification [1] to %s: %s", Recipient, err.Error()))
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification [2] to %s: %s", Recipient, err.Error()))
		return err
	}

	_, err = w.Write([]byte(completeMessage))
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification [3] to %s: %s", Recipient, err.Error()))
		return err
	}

	err = w.Close()
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Can't send E-Mail notification [4] to %s: %s", Recipient, err.Error()))
		return err
	}

	c.Quit()

	return nil
}
