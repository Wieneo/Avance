package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetQueue returns the queue struct to the given id
func GetQueue(QueueID int64, Project int64) (models.Queue, bool, error) {
	var queue models.Queue
	err := Connection.QueryRow(`SELECT "ID", "Name" FROM "Queue" WHERE "ID" = $1 AND "Project" = $2`, QueueID, Project).Scan(&queue.ID, &queue.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Queue{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Queue{}, true, err
	}

	return queue, true, nil
}

//GetQueueUNSAFE returns the queue struct to the given id ignoring the project relationship
func GetQueueUNSAFE(QueueID int64) (models.Queue, bool, error) {
	var queue models.Queue
	err := Connection.QueryRow(`SELECT "ID", "Name" FROM "Queue" WHERE "ID" = $1`, QueueID).Scan(&queue.ID, &queue.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Queue{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Queue{}, true, err
	}

	return queue, true, nil
}

//QueuesInProject returns all queues from a project
func QueuesInProject(Project models.Project) ([]models.Queue, error) {
	Queues := make([]models.Queue, 0)
	rows, err := Connection.Query(`SELECT "ID", "Name" FROM "Queue" WHERE "Project" = $1`, Project.ID)
	if err != nil {
		return make([]models.Queue, 0), err
	}

	for rows.Next() {
		var SingleQueue models.Queue
		rows.Scan(&SingleQueue.ID, &SingleQueue.Name)
		Queues = append(Queues, SingleQueue)
	}

	rows.Close()

	return Queues, nil
}

//CreateQueue creates a queue in the database
func CreateQueue(Name string, Project int64) (int64, error) {
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Queue" ("Name", "Project") VALUES ($1, $2) RETURNING "ID"`, Name, Project).Scan(&newID)
	if err != nil {
		return 0, err
	}

	return newID, nil
}

//PatchQueue patches the given Queue
func PatchQueue(Queue models.Queue) error {
	_, err := Connection.Exec(`UPDATE "Queue" SET "Name" = $1 WHERE "ID" = $2`, Queue.Name, Queue.ID)
	return err
}

//RemoveQueue removes a queue from the database
func RemoveQueue(Project int64, Queue int64) error {
	//I know that project isn't really needed as queue ids are unique anyway
	//Its just a safety measure ;)
	_, err := Connection.Exec(`DELETE FROM "Queue" WHERE "ID" = $1 AND "Project" = $2`, Queue, Project)
	return err
}
