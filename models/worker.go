package models

import (
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
	ID       int64
	Type     WorkerTaskType
	Data     string
	QueuedAt time.Time
	Status   WorkerTaskStatus
}

//WorkerTaskType stores the type of the worker task
type WorkerTaskType int

const (
	//DeleteUser triggers the user deletion
	DeleteUser WorkerTaskType = iota
)

func (e WorkerTaskType) String() string {
	switch e {
	case DeleteUser:
		return "Delete User"
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
