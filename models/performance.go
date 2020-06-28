package models

//WantedProperties is used to tell the GetTicket function what properties to load
type WantedProperties struct {
	All, Owner, Queue, Severity, Status, Relations, Recipients, Actions bool
}
