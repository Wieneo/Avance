package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetQueue returns the queue struct to the given id
func GetQueue(ProjectID int64, QueueID int64) (models.Queue, bool, error) {
	var queue models.Queue
	err := Connection.QueryRow(`SELECT "ID", "Name" FROM "Queue" WHERE "ID" = $1 AND "Project" = $2`, QueueID, ProjectID).Scan(&queue.ID, &queue.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return models.Queue{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Queue{}, true, err
	}

	return queue, true, nil
}

//GetQueueUNSAFE returns the queue struct to the given id without checking if its contained in a project
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
