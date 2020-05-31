package db

import (
	"fmt"
	"os"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetAllWorkers returns all workers from the database
func GetAllWorkers() ([]models.Worker, error) {
	allworkers := make([]models.Worker, 0)
	rows, err := Connection.Query(`SELECT "ID", "Name", "LastSeen", "Active" FROM "Workers"`)
	if err != nil {
		return allworkers, err
	}

	for rows.Next() {
		var singleWorker models.Worker
		rows.Scan(&singleWorker.ID, &singleWorker.Name, &singleWorker.LastSeen, &singleWorker.Active)
		allworkers = append(allworkers, singleWorker)
	}

	return allworkers, nil
}

//RegisterWorker tries to find out if the worker already is registered and if not registers it
func RegisterWorker() (int, error) {
	name, _ := os.Hostname()

	rows, err := Connection.Query(`SELECT "ID" FROM "Workers" WHERE "Name" = $1`, name)
	if err != nil {
		return 0, err
	}

	var wid int
	if rows.Next() {
		rows.Scan(&wid)

		dev.LogInfo(fmt.Sprintf("Worker already was registered: %d", wid))

		return wid, nil
	}

	rows.Close()

	//Register
	dev.LogInfo("Worker seems to be new! Registering...")
	err = Connection.QueryRow(`INSERT INTO "Workers" ("Name", "LastSeen", "Active") VALUES ($1, $2, $3) RETURNING "ID"`, name, time.Now(), true).Scan(&wid)
	if err != nil {
		return wid, err
	}

	dev.LogInfo(fmt.Sprintf("New worker registered with ID %d", wid))
	return wid, nil
}

//GetWorkerStatus returns wheter the worker is active or not
func GetWorkerStatus() (bool, error) {
	status := false
	name, _ := os.Hostname()
	err := Connection.QueryRow(`SELECT "Active" FROM "Workers" WHERE "Name" = $1`, name).Scan(&status)

	if err := refreshWorker(); err != nil {
		dev.LogError(err, "Couldn't update workers 'LastSeen' Date: "+err.Error())
	}

	return status, err
}

func refreshWorker() error {
	name, _ := os.Hostname()
	_, err := Connection.Exec(`UPDATE "Workers" SET "LastSeen" = $1 WHERE "Name" = $2`, time.Now(), name)
	return err
}
