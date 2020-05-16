package db

import (
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
)

//GetQueue returns the queue struct to the given id
func GetQueue(QueueID int64) (models.Queue, error) {
	var queue models.Queue
	err := Connection.QueryRow(`SELECT "ID", "Name" FROM "Queue" WHERE "ID" = $1`, QueueID).Scan(&queue.ID, &queue.Name)
	if err != nil {
		dev.LogError(err, err.Error())
		return models.Queue{}, err
	}

	return queue, nil
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
