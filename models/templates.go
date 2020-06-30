package models

import (
	"database/sql"
	"fmt"
	"time"
)

//Constants

//DateFormat is used as the default format for parsing dates the user specifies
const DateFormat = "2006-01-02T15:04:05.000Z"

//GetAllowedImageFormates defines what Image formates a Prfile Picture is allowed to be
func GetAllowedImageFormates() []string {
	return []string{"png", "jpg", "jpeg", "gif"}
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
	Recipients   RecipientCollection
}

//CreateTicket is used for passing all necessary information to the db subroutine for creating tickets
type CreateTicket struct {
	Title         string
	Description   string
	Queue         int64
	OwnedByNobody bool
	Owner         int64
	Severity      Severity
	Status        Status
	IsStalled     bool
	StalledUntil  string
}

//RecipientCollection stores all recipients assigned to a single ticket
type RecipientCollection struct {
	Requestors, Readers, Admins []Recipient
}

//Recipient stores a single recipient (mail / known user)
type Recipient struct {
	ID int64
	//RecipientType is only populated when used with AllRecipients
	Type RecipientType
	User struct {
		Valid bool
		Value User
	}
	Mail string
}

//RecipientType is only used to map entries from the database into the 3 arrays used in the RecipientCollection struct
type RecipientType int

const (
	//Requestors -> See RecipientCollection
	Requestors RecipientType = iota
	//Readers -> See RecipientCollection
	Readers
	//Admins -> See RecipientCollection
	Admins
)

func (e RecipientType) String() string {
	switch e {
	case Requestors:
		return "Requestors"
	case Readers:
		return "Readers"
	case Admins:
		return "Admins"
	default:
		return fmt.Sprintf("%d", int(e))
	}
}

//AllRecipients returns all recipients from admins, requestors, readers
func (e Ticket) AllRecipients() []Recipient {
	allrecp := make([]Recipient, 0)
	for _, k := range e.Recipients.Admins {
		k.Type = Admins
		allrecp = append(allrecp, k)
	}
	for _, k := range e.Recipients.Requestors {
		k.Type = Requestors
		allrecp = append(allrecp, k)
	}
	for _, k := range e.Recipients.Readers {
		k.Type = Readers
		allrecp = append(allrecp, k)
	}
	return allrecp
}

//Action stores a single update about a ticket
type Action struct {
	ID       int64
	Type     ActionType
	Title    string
	Content  string
	IssuedAt time.Time
	IssuedBy Issuer
	Tasks    []int64
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
