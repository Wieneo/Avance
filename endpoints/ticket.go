package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
)

//GetTicket returns the serialized ticket to the user
func GetTicket(w http.ResponseWriter, r *http.Request) {
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)
	ticket, found, err := db.GetTicket(ticketid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Ticket not found")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, ticket.Queue)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if allperms.Admin || perms.CanSee {
		json.NewEncoder(w).Encode(ticket)
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You don't have access to that queue!")
	}
}

//GetTicketsFromQueue returns all tickets in a given queue
func GetTicketsFromQueue(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Project/Queue not found")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		json.NewEncoder(w).Encode(tickets)
	} else {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You don't have access to that queue!")
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

//CreateTicketsInQueue returns all tickets in a given queue
func CreateTicketsInQueue(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)

	queue, found, err := db.GetQueue(projectid, queueid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Project/Queue not found")
		return
	}

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	allperms, perms, err := perms.GetPermissionsToQueue(user, queue)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
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
			w.WriteHeader(500)
			dev.ReportError(err, w, err.Error())
			return
		}

		if !found {
			w.WriteHeader(404)
			dev.ReportUserError(w, "Owner "+req.Owner+" couldn't be found")
			return
		}

		ownedByNobody = false
	}

	statusid, found, err := db.SearchStatus(req.Status)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Status "+req.Status+" couldn't be found")
		return
	}

	severityid, found, err := db.SearchSeverity(req.Severity)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Severity "+req.Severity+" couldn't be found")
		return
	}

	isStalled := false

	if !utils.IsEmpty(req.StalledUntil) {
		_, err := time.Parse("2006-01-02T15:04:05.000Z", req.StalledUntil)
		if err != nil {
			w.WriteHeader(406)
			dev.ReportUserError(w, "StalledUntil isn't a valid date/time")
			return
		}

		isStalled = true
	}

	//Now everything should be ok
	id, err := db.CreateTicket(req.Title, req.Description, queueid, ownedByNobody, ownerID, severityid, statusid, isStalled, req.StalledUntil)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	ticket, found, err := db.GetTicket(id)
	if err != nil || !found {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(ticket)
}
