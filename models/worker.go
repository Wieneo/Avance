package models

import (
	"database/sql"
	"fmt"
	"time"
)

//Worker represents a single worker from the database
type Worker struct {
	ID       int
	Name     string
	LastSeen time.Time
	Active   bool
}

//WorkerTask is the generic task for a worker
type WorkerTask struct {
	ID        int64
	Type      WorkerTaskType
	Data      string
	QueuedAt  time.Time
	Status    WorkerTaskStatus
	Interval  sql.NullInt32
	LastRun   sql.NullTime
	Recipient sql.NullString
	Ticket    sql.NullInt64
}

//NotificationCollection contains all collected notifications about a ticket
type NotificationCollection struct {
	Subject       string
	Notifications []Notification
}

//Notification stores a single notificaiton
type Notification struct {
	Title, Content string
}

//WorkerTaskType stores the type of the worker task
type WorkerTaskType int

const (
	//DeleteUser triggers the user deletion
	DeleteUser WorkerTaskType = iota
	//Debug is used during development to test worker behaviour
	Debug
	//SendNotification informs the user about ticket updates
	SendNotification
)

func (e WorkerTaskType) String() string {
	switch e {
	case DeleteUser:
		return "Delete User"
	case SendNotification:
		return "Send Notification"
	case Debug:
		return "Debug"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

//WorkerTaskStatus stores the status of the worker task
type WorkerTaskStatus int

const (
	//Idle is set by default to signal the job is ready to be picked up
	Idle WorkerTaskStatus = iota
	//InProgress is set if it got picked up by a worker
	InProgress
	//Failed is set if job couldn't execute
	Failed
	//Finished is set after job completion
	Finished
)
