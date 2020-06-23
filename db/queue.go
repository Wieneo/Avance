package db

import (
	"fmt"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
)

//GetQueue returns the queue struct to the given id
func GetQueue(ProjectID, QueueID int64) (models.Queue, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting queue %d in project %d", QueueID, ProjectID))
	var queue models.Queue
	err := Connection.QueryRow(`SELECT "ID", "Name" FROM "Queue" WHERE "ID" = $1 AND "Project" = $2`, QueueID, ProjectID).Scan(&queue.ID, &queue.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			dev.LogDebug(fmt.Sprintf("[DB] Requested queue (Q: %d, P: %d) wasnt found -> Returning empty queue struct", QueueID, ProjectID))
			return models.Queue{}, false, nil
		}

		dev.LogError(err, err.Error())
		return models.Queue{}, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Got queue %d: Name: %s", QueueID, queue.Name))
	return queue, true, nil
}

//GetQueueUNSAFE returns the queue struct to the given id without checking if its contained in a project
func GetQueueUNSAFE(QueueID int64) (models.Queue, bool, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting queue %d in UNSAFE mode", QueueID))
	var queue models.Queue
	err := Connection.QueryRow(`SELECT "ID", "Name" FROM "Queue" WHERE "ID" = $1`, QueueID).Scan(&queue.ID, &queue.Name)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			dev.LogDebug(fmt.Sprintf("[DB] Requested queue %d wasnt found -> Returning empty queue struct", QueueID))
			return models.Queue{}, false, nil
		}
		dev.LogError(err, err.Error())
		return models.Queue{}, true, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Got queue %d: Name: %s", QueueID, queue.Name))

	return queue, true, nil
}

//QueuesInProject returns all queues from a project
func QueuesInProject(Project models.Project) ([]models.Queue, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Getting all queues in project %d", Project.ID))
	Queues := make([]models.Queue, 0)
	rows, err := Connection.Query(`SELECT "ID", "Name" FROM "Queue" WHERE "Project" = $1`, Project.ID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] Error happened while getting queues for project %d: %s", Project.ID, err.Error()))
		return make([]models.Queue, 0), err
	}

	for rows.Next() {
		var SingleQueue models.Queue
		rows.Scan(&SingleQueue.ID, &SingleQueue.Name)
		Queues = append(Queues, SingleQueue)
	}

	rows.Close()

	dev.LogDebug(fmt.Sprintf("[DB] Got %d queues for project %d", len(Queues), Project.ID))

	return Queues, nil
}

//CreateQueue creates a queue in the database
func CreateQueue(Name string, Project int64) (int64, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Creating Queue '%s' in project %d", Name, Project))
	var newID int64
	err := Connection.QueryRow(`INSERT INTO "Queue" ("Name", "Project") VALUES ($1, $2) RETURNING "ID"`, Name, Project).Scan(&newID)
	if err != nil {
		return 0, err
	}

	dev.LogDebug(fmt.Sprintf("[DB] Queue with id %d created", newID))
	return newID, nil
}

//PatchQueue patches the given Queue
func PatchQueue(Queue models.Queue) error {
	dev.LogDebug(fmt.Sprintf("[DB] Patching queue %d", Queue.ID))
	_, err := Connection.Exec(`UPDATE "Queue" SET "Name" = $1 WHERE "ID" = $2`, Queue.Name, Queue.ID)
	dev.LogDebug(fmt.Sprintf("[DB] Patch completed for queue %d", Queue.ID))
	return err
}

//RemoveQueue removes a queue from the database
func RemoveQueue(Project, Queue int64) error {
	//I know that project isn't really needed as queue ids are unique anyway
	//Its just a safety measure ;)
	dev.LogDebug(fmt.Sprintf("[DB] Removing queue %d from project %d", Queue, Project))
	_, err := Connection.Exec(`DELETE FROM "Queue" WHERE "ID" = $1 AND "Project" = $2`, Queue, Project)
	return err
}

//GetProjectFromQueue returns the project the queue is assigned to
func GetProjectFromQueue(QueueID int64) (models.Project, error) {
	dev.LogDebug(fmt.Sprintf("[DB] Reverse lookup -> Getting project for queue %d", QueueID))
	var projectID int64
	err := Connection.QueryRow(`SELECT "Project" FROM "Queue" WHERE "ID" = $1`, QueueID).Scan(&projectID)
	if err != nil {
		dev.LogDebug(fmt.Sprintf("[DB] An error happened while getting project for queue %d: %s", QueueID, err.Error()))
		return models.Project{}, err
	}

	project, _, err := GetProject(projectID)
	dev.LogDebug(fmt.Sprintf("[DB] Reverse lookup for queue %d returned: ID: %d, Name: %s", QueueID, project.ID, project.Name))
	return project, err
}
