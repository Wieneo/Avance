package models

import (
	"database/sql"
	"time"
)

///Constants

//GetAllowedImageFormates defines what Image formates a Prfile Picture is allowed to be
func GetAllowedImageFormates() []string {
	return []string{"png", "jpg", "jpeg", "gif"}
}

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
	ID   int64
	Name string
}

//Permissions should store all permissions regarding a user or a group
type Permissions struct {
	Admin                      bool
	CanCreateUsers             bool
	CanModifyUsers             bool
	CanDeleteUsers             bool
	CanCreateGroups            bool
	CanModifyGroups            bool
	CanDeleteGroups            bool
	CanChangePermissionsGlobal bool
	CanSeeWorker               bool
	CanChangeWorker            bool
	AccessTo                   struct {
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
	Relations    []Relation
	Actions      []Action
}

//Action stores a single update about a ticket
type Action struct {
	ID       int64
	Type     ActionType
	Title    string
	Content  string
	IssuedAt time.Time
	IssuedBy Issuer
}

//Issuer stores wheter a action has a valid issuer and which one
type Issuer struct {
	Valid  bool
	Issuer User
}

//ActionType defines to type of an action to be a comment / answer / etc...
type ActionType int

const (
	//Answer is set when a answer is issued
	Answer ActionType = iota
	//Comment is set when a comment is issued
	Comment
	//PropertyUpdate gets set if properties of the ticket got changed
	PropertyUpdate
	//Unspecific is a generic type
	Unspecific
)

//RelationType stores the type of relation between tickets
type RelationType int

const (
	//References is set if a ticket references another ticket
	References RelationType = iota
	//ReferencedBy is set if another ticket references the ticket
	ReferencedBy
	//ParentOf is set if the ticket is a parent
	ParentOf
	//ChildOf is set if the ticket is a child
	ChildOf
)

//Relation stores the relation between two tickets
//Should be stored like "Ticket which contains the relation $TYPE the other ticket from the struct"
type Relation struct {
	ID          int64
	OtherTicket Ticket
	Type        RelationType
}
