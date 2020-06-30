package endpoints

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//GetTicket returns the requested ticket
func GetTicket(w http.ResponseWriter, r *http.Request) {
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

	ticket, found, err := db.GetTicketUnsafe(ticketid, models.WantedProperties{All: true})
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.TicketNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	project, err := db.GetProjectFromQueue(ticket.Queue.ID)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, pperms, err := perms.GetPermissionsToProject(user, project)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !allperms.Admin {
		if !pperms.CanSee {
			w.WriteHeader(403)
			dev.ReportUserError(w, "You don't have access to that project!")
			return
		}

		_, qperms, err := perms.GetPermissionsToQueue(user, ticket.Queue)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !qperms.CanSee {
			w.WriteHeader(403)
			dev.ReportUserError(w, templates.QueueNoPerms)
			return
		}
	}

	json.NewEncoder(w).Encode(ticket)
}

//GetTicketFullPath returns the serialized ticket to the user
func GetTicketFullPath(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if allperms.Admin || perms.CanSee {
		ticket, found, err := db.GetTicket(ticketid, queueid, models.WantedProperties{All: true})
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, templates.TicketNotFound)
			return
		}

		json.NewEncoder(w).Encode(ticket)
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, templates.QueueNoPerms)
	}
}

//GetTicketsFromQueue returns all tickets in a given queue
func GetTicketsFromQueue(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if allperms.Admin || perms.CanSee {
		var showInvisible bool

		if len(r.URL.Query()["showInvisible"]) > 0 {
			showInvisible, err = strconv.ParseBool(r.URL.Query()["showInvisible"][0])
			if err != nil {
				w.WriteHeader(400)
				dev.ReportUserError(w, "showInvisible Argument is not a boolean")
				return
			}
		}

		tickets, err := db.GetTicketsInQueue(queueid, showInvisible)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		json.NewEncoder(w).Encode(tickets)
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, templates.QueueNoPerms)
	}
}

type ticketWebRequest struct {
	Title        string
	Description  string
	Owner        string //Username instead of ID
	Severity     string //Name instead of ID
	Status       string //Name instead of ID
	StalledUntil string
}

//CreateTicketsInQueue creates a ticket
func CreateTicketsInQueue(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !allperms.Admin && !perms.CanCreateTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to create tickets in that queue.")
		return
	}

	var req ticketWebRequest

	rawBytes, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed")
		return
	}

	//Check if all required fields are filled
	if utils.IsEmpty(req.Title) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Title can't be empty")
		return
	}

	if utils.IsEmpty(req.Status) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Status can't be empty")
		return
	}

	if utils.IsEmpty(req.Severity) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Severity can't be empty")
		return
	}

	var ownerID int64
	ownedByNobody := true

	//If owner is not specified -> Nobody is the owner
	//This is checked again in the database part
	if !utils.IsEmpty(req.Owner) {
		ownerID, found, err = db.SearchUser(req.Owner)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Owner "+req.Owner+" couldn't be found")
			return
		}

		ownedByNobody = false
	}

	statusid, found, err := db.SearchStatus(projectid, req.Status)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Status "+req.Status+" couldn't be found")
		return
	}

	severityid, found, err := db.SearchSeverity(projectid, req.Severity)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	severity, _, err := db.GetSeverity(projectid, severityid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	status, _, err := db.GetStatus(projectid, statusid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Severity "+req.Severity+" couldn't be found")
		return
	}

	isStalled := false

	if !utils.IsEmpty(req.StalledUntil) {
		_, err := time.Parse(models.DateFormat, req.StalledUntil)
		if err != nil {
			w.WriteHeader(406)
			dev.ReportUserError(w, "StalledUntil isn't a valid date/time")
			return
		}

		isStalled = true
	}

	newTicket := models.CreateTicket{
		Title:         req.Title,
		Description:   req.Description,
		Queue:         queueid,
		OwnedByNobody: ownedByNobody,
		Owner:         ownerID,
		Severity:      severity,
		Status:        status,
		IsStalled:     isStalled,
		StalledUntil:  req.StalledUntil,
	}

	//Now everything should be ok
	id, err := db.CreateTicket(newTicket)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	ticket, found, err := db.GetTicket(id, queueid, models.WantedProperties{All: true})
	if err != nil || !found {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

//PatchTicketsInQueue patches a ticket
func PatchTicketsInQueue(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !allperms.Admin && !perms.CanEditTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to patch tickets in that queue.")
		return
	}

	var req ticketWebRequest

	rawBytes, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawBytes, &req)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed")
		return
	}

	ticket, found, err := db.GetTicket(ticketid, queueid, models.WantedProperties{All: true})
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.TicketNotFound)
		return
	}

	somethingChanged := false

	if !utils.IsEmpty(req.Title) && req.Title != ticket.Title {
		ticket.Title = req.Title
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Description) && req.Description != ticket.Description {
		ticket.Description = req.Description
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Owner) && req.Owner != ticket.Owner.Username {
		ownerID, found, err := db.SearchUser(req.Owner)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Owner "+req.Owner+" couldn't be found")
			return
		}

		ticket.OwnerID.Valid = true
		ticket.OwnerID.Int64 = ownerID
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Status) && req.Status != ticket.Status.Name {
		statusid, found, err := db.SearchStatus(projectid, req.Status)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Status "+req.Status+" couldn't be found")
			return
		}

		ticket.StatusID = statusid
		somethingChanged = true
	}

	if !utils.IsEmpty(req.Severity) && req.Severity != ticket.Severity.Name {
		severityid, found, err := db.SearchSeverity(projectid, req.Severity)
		if err != nil {

			utils.ReportInternalErrorToUser(err, w)
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Severity "+req.Severity+" couldn't be found")
			return
		}

		ticket.SeverityID = severityid
		somethingChanged = true
	}

	if !utils.IsEmpty(req.StalledUntil) && req.StalledUntil != ticket.StalledUntil.Time.Format(models.DateFormat) {
		t, err := time.Parse(models.DateFormat, req.StalledUntil)
		if err != nil {
			w.WriteHeader(406)
			dev.ReportUserError(w, "StalledUntil isn't a valid date/time")
			return
		}

		ticket.StalledUntil.Time = t
		ticket.StalledUntil.Valid = true
		somethingChanged = true
	}

	if !somethingChanged {
		w.WriteHeader(406)
		dev.ReportUserError(w, templates.NothingChanged)
		return
	}

	ticket, err = db.PatchTicket(ticket)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

//DeletePropertyFromTicket removes owner / stalleduntil from ticket
func DeletePropertyFromTicket(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.QueueNotFound)
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !allperms.Admin && !perms.CanEditTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to patch tickets in that queue.")
		return
	}

	ticket, found, err := db.GetTicket(ticketid, queueid, models.WantedProperties{All: true})
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.TicketNotFound)
		return
	}

	if strings.Split(r.RequestURI, "/")[9] == "owner" {
		if !ticket.OwnerID.Valid {
			w.WriteHeader(406)
			dev.ReportUserError(w, templates.NothingChanged)
			return
		}

		ticket.OwnerID.Valid = false
		ticket.OwnerID.Int64 = 0

	} else if strings.Split(r.RequestURI, "/")[9] == "stalleduntil" {
		if !ticket.StalledUntil.Valid {
			w.WriteHeader(406)
			dev.ReportUserError(w, templates.NothingChanged)
			return
		}

		var newStalled sql.NullTime
		ticket.StalledUntil = newStalled
	}

	ticket, err = db.PatchTicket(ticket)
	if err != nil {

		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(ticket)
}

type recipientRequest struct {
	Type models.RecipientType
	User string
	Mail string
}
