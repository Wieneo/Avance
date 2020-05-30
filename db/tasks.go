package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//ReserveTask reserves a task for this worker
func ReserveTask() (int64, error) {
	var newID int64
	err := Connection.QueryRow(`SELECT public."GetTask"();`).Scan(&newID)

	return newID, err
}

//GetTask returns the task to a give ID
func GetTask(TaskID int64) (models.WorkerTask, error) {
	var workerTask models.WorkerTask
	err := Connection.QueryRow(`SELECT "ID", "Task", "QueuedAt", "Status", "Type" FROM "Tasks" WHERE "ID" = $1`, TaskID).Scan(&workerTask.ID, &workerTask.Data, &workerTask.QueuedAt, &workerTask.Status, &workerTask.Type)
	if err != nil {
		return models.WorkerTask{}, err
	}

	return workerTask, err
}

//PatchTask patches the task in the database. Ignores QueuedAt!
func PatchTask(Task models.WorkerTask) error {
	_, err := Connection.Exec(`UPDATE "Tasks" SET "Task" = $1, "Status" = $2, "Type" = $3 WHERE "ID" = $4`, Task.Data, Task.Status, Task.Type, Task.ID)
	return err
}
