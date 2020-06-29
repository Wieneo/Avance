package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//ReserveTask reserves a task for this worker
func ReserveTask() (models.WorkerTask, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Reserving task via public.'GetTask'()"))
	var newID int64
	err := Connection.QueryRow(`SELECT public."GetTask"();`).Scan(&newID)

	if err != nil {
		return models.WorkerTask{}, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Got task id %d", newID))

	dev.LogDebug(fmt.Sprintf("[DB] Looking up task %d", newID))
	task, _, err := GetTask(newID)

	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error while looking up task %d -> Returning empty task struct: %s", newID, err.Error()))
		return models.WorkerTask{}, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Updating task %d to reflect pickup", newID))
	hostname, _ := os.Hostname()
	err = AddResult(&task, "Job picked up by worker "+hostname)

	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error while adding result to task %d -> Returning empty task struct: %s", newID, err.Error()))
		return models.WorkerTask{}, err
	}

	return task, err
}

//CreateTask queues a task
func CreateTask(Type models.WorkerTaskType, Data string, Interval sql.NullInt32, Recipient sql.NullString, Ticket sql.NullInt64, NotificationType sql.NullInt32) (int64, error) {
	var taskID int64
	err := Connection.QueryRow(`INSERT INTO "Tasks" ("Task", "QueuedAt", "Status", "Type", "Interval", "Recipient", "Ticket", "LastRun", "NotifiationType") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "ID"`, Data, time.Now(), models.Idle, Type, Interval, Recipient, Ticket, time.Now(), NotificationType).Scan(&taskID)
	return taskID, err
}

//GetTask returns the task to a give ID
func GetTask(TaskID int64) (models.WorkerTask, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting task %d", TaskID))

	var rawResults sql.NullString
	var workerTask models.WorkerTask
	err := Connection.QueryRow(`SELECT "ID", "Task", "QueuedAt", "Status", "Type", "Interval", "LastRun", "Recipient", "Ticket", "Results" FROM "Tasks" WHERE "ID" = $1`, TaskID).Scan(&workerTask.ID, &workerTask.Data, &workerTask.QueuedAt, &workerTask.Status, &workerTask.Type, &workerTask.Interval, &workerTask.LastRun, &workerTask.Recipient, &workerTask.Ticket, &rawResults)
	if err != nil {
		if err == sql.ErrNoRows {
			dev.LogDebug(fmt.Sprintf("[DB] Requested task %d wasnt found -> Returning empty queue struct", TaskID))
			return models.WorkerTask{}, false, nil
		}

		dev.LogDebug(fmt.Sprintf("[DB] Error happened while retrieving task %d -> Returning empty task struct: %s", TaskID, err.Error()))
		return models.WorkerTask{}, true, err
	}

	if rawResults.Valid {
		err = json.Unmarshal([]byte(rawResults.String), &workerTask.Results)
		if err != nil {
			dev.LogDebug(fmt.Sprintf("[DB] Error happened while parsing results of task %d -> Returning empty task struct: %s", TaskID, err.Error()))
			return models.WorkerTask{}, true, err
		}
	} else {
		workerTask.Results = make([]models.TaskResult, 0)
	}

	dev.LogDebug(fmt.Sprintf("[DB] Got task %d", workerTask.ID))
	return workerTask, true, err
}

//PatchTask patches the task in the database. Ignores QueuedAt!
func PatchTask(Task models.WorkerTask) error {
	dev.LogDebug(fmt.Sprintf("[DB] Patching task %d", Task.ID))
	_, err := Connection.Exec(`UPDATE "Tasks" SET "Task" = $1, "Status" = $2, "Type" = $3, "Interval" = $4, "LastRun" = $5 WHERE "ID" = $6`, Task.Data, Task.Status, Task.Type, Task.Interval, Task.LastRun, Task.ID)
	if err == nil {
		dev.LogDebug(fmt.Sprintf("[DB] Patched task %d", Task.ID))
	}
	return err
}

//AddResult adds a result to the results array of the task
func AddResult(Task *models.WorkerTask, Result string) error {
	dev.LogDebug(fmt.Sprintf("[DB] Adding result to task %d", Task.ID))
	Task.Results = append(Task.Results, models.TaskResult{IssuedAt: time.Now(), Result: Result})

	rawJSON, _ := json.Marshal(Task.Results)
	_, err := Connection.Exec(`UPDATE "Tasks" SET "Results" = $1 WHERE "ID" = $2`, rawJSON, Task.ID)

	dev.LogDebug(fmt.Sprintf("[DB] Added result to task %d", Task.ID))
	return err
}
