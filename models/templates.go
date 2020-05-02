package models

//User is the default user struct from the database
type User struct {
	ID          int
	Username    string
	Mail        string
	Permissions Permissions
}

//Group is the default Group struct from the database
type Group struct {
	ID          int
	Name        string
	Permissions Permissions
}

//Project is the default Project struct from the database
type Project struct {
	ID   int
	Name string
}

//Queue is the default Queue struct from the database
type Queue struct {
	ID      int
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

//ProjectPermission stores the permissions given to a single project
type ProjectPermission struct {
	ProjectID            int
	CanSee               bool
	CanModify            bool
	CanChangePermissions bool
	CanCreateQueues      bool
	CanModifyQueues      bool
	CanRemoveQueues      bool
}

//QueuePermission stores the permissions given to a single queue
type QueuePermission struct {
	QueueID              int
	CanSee               bool
	CanCreateTicket      bool
	CanEditTicket        bool
	CanModify            bool
	CanChangePermissions bool
}
