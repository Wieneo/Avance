package worker

import (
	"errors"
	"os"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/worker/smtp"
)

//InitWorker initializes all parts of the application so the worker can start operating
func InitWorker() {
	dev.LogInfo("Starting Worker Version", utils.WorkerVersion, utils.WorkerChannel)
	hostname, _ := os.Hostname()
	dev.LogInfo("Hostname: " + hostname)
	config.LoadConfig()
	db.Init(false)

	_, err := db.RegisterWorker()
	if err != nil {
		dev.LogFatal(err, "Couldn't register worker: "+err.Error())
	}

	//Check if everything is working
	state := getHealthState()
	if len(state.Errors) > 0 {
		dev.LogError(errors.New("Startup Error"), "The following errors occured during startup: ")
		for _, k := range state.Errors {
			dev.LogError(errors.New("Startup Error"), k)
		}
		os.Exit(1)
	}

	dev.LogInfo("Everything looks fine. Let's get to work!")

	StartQueueService()
	if config.CurrentConfig.Worker.Listen {
		StartServing()
	} else {
		select {}
	}
}

func getHealthState() models.WorkerHealth {
	errors := make([]string, 0)

	var dummyDBVersion string

	dBAlive := true
	smtpAlive := true

	//Check Database connection
	err := db.Connection.QueryRow(`SELECT "Name" FROM "Patches" LIMIT 1`).Scan(&dummyDBVersion)
	if err != nil {
		dBAlive = false
		errors = append(errors, err.Error())
	}

	//Check SMTP Connection
	if config.CurrentConfig.SMTP.Enabled {
		if c, err := smtp.Login(); err != nil {
			smtpAlive = false
			errors = append(errors, err.Error())
		} else {
			c.Quit()
		}
	}

	return models.WorkerHealth{
		DBAlive:   dBAlive,
		SMTPAlive: smtpAlive,
		Errors:    errors,
	}
}
