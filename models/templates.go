package models

import (
	"database/sql"
	"time"
)

//User is the default user struct from the database
type User struct {
	ID          int64
	Username    string
	Mail        string
	Permissions Permissions
	Firstname   string
	Lastname    string
}

//Group is the default Group struct from the database
type Group struct {
	ID          int64
	Name        string
	Permissions Permissions
}

//Project is the default Project struct from the database
type Project struct {
	ID          int64
	Name        string
	Description string
}

//Queue is the default Queue struct from the database
type Queue struct {
	ID      int64
	Name    string
	Project Project
}

//Permissions should store all permissions regarding a user or a group
type Permissions struct {
	Admin    bool
	AccessTo struct {
		Projects []ProjectPermission
		Queues   []QueuePermission
	}
}

//Severity is the default Severity struct from the database
type Severity struct {
	ID           int64
	Enabled      bool
	Name         string
	DisplayColor string
	Priority     int
}

//Status is the default Status struct from the datbabase
type Status struct {
	ID             int64
	Enabled        bool
	Name           string
	DisplayColor   string
	TicketsVisible bool
}

/*
	ALWAYS!! If a new permission is added here! PLEASE add it to perms/combine.go
*/

//ProjectPermission stores the permissions given to a single project
type ProjectPermission struct {
	ProjectID            int64
	CanSee               bool
	CanModify            bool
	CanChangePermissions bool
	CanCreateQueues      bool
	CanModifyQueues      bool
	CanRemoveQueues      bool
	CanCreateSeverities  bool
	CanModifySeverities  bool
	CanRemoveSeverities  bool
	CanCreateStatuses    bool
	CanModifyStatuses    bool
	CanRemoveStatuses    bool
}

//QueuePermission stores the permissions given to a single queue
type QueuePermission struct {
	QueueID              int64
	CanSee               bool
	CanCreateTicket      bool
	CanEditTicket        bool
	CanModify            bool
	CanChangePermissions bool
}

//Ticket stores all information about a ticket
type Ticket struct {
	ID           int64
	Title        string
	Description  string
	QueueID      int64
	Queue        Queue
	OwnerID      sql.NullInt64
	Owner        User
	SeverityID   int64
	Severity     Severity
	StatusID     int64
	Status       Status
	CreatedAt    time.Time
	LastModified time.Time
	StalledUntil sql.NullTime
	Meta         string //Needs to be changed later!
}
