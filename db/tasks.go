package db

import (
	"database/sql"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//ReserveTask reserves a task for this worker
func ReserveTask() (int64, error) {
	var newID int64
	err := Connection.QueryRow(`SELECT public."GetTask"();`).Scan(&newID)

	return newID, err
}

//CreateTask queues a task
func CreateTask(Type models.WorkerTaskType, Data string, Interval sql.NullInt32, Recipient sql.NullString, Ticket sql.NullInt64, NotificationType sql.NullInt32) (int64, error) {
	var taskID int64
	err := Connection.QueryRow(`INSERT INTO "Tasks" ("Task", "QueuedAt", "Status", "Type", "Interval", "Recipient", "Ticket", "LastRun", "NotifiationType") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING "ID"`, Data, time.Now(), models.Idle, Type, Interval, Recipient, Ticket, time.Now(), NotificationType).Scan(&taskID)
	return taskID, err
}

//GetTask returns the task to a give ID
func GetTask(TaskID int64) (models.WorkerTask, error) {
	var workerTask models.WorkerTask
	err := Connection.QueryRow(`SELECT "ID", "Task", "QueuedAt", "Status", "Type", "Interval", "LastRun", "Recipient", "Ticket" FROM "Tasks" WHERE "ID" = $1`, TaskID).Scan(&workerTask.ID, &workerTask.Data, &workerTask.QueuedAt, &workerTask.Status, &workerTask.Type, &workerTask.Interval, &workerTask.LastRun, &workerTask.Recipient, &workerTask.Ticket)
	if err != nil {
		return models.WorkerTask{}, err
	}

	return workerTask, err
}

//PatchTask patches the task in the database. Ignores QueuedAt!
func PatchTask(Task models.WorkerTask) error {
	_, err := Connection.Exec(`UPDATE "Tasks" SET "Task" = $1, "Status" = $2, "Type" = $3, "Interval" = $4, "LastRun" = $5 WHERE "ID" = $6`, Task.Data, Task.Status, Task.Type, Task.Interval, Task.LastRun, Task.ID)
	return err
}
