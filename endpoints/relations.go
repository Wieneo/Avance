package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

type relationWebRequest struct {
	OtherTicket  int64
	RelationType models.RelationType
}

//CreateRelation adds a relation to a ticket
func CreateRelation(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)

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

	if !allperms.Admin && !perms.CanEditTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to patch tickets in that queue.")
		return
	}

	ticket, found, err := db.GetTicket(ticketid, queueid, true)
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

	var relation relationWebRequest

	rawBytes, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawBytes, &relation)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed")
		return
	}

	_, found, err = db.GetTicket(relation.OtherTicket, queueid, false)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "OtherTicket not found")
		return
	}

	switch relation.RelationType {
	case models.References, models.ReferencedBy, models.ParentOf, models.ChildOf:
		{
			break
		}
	default:
		{
			w.WriteHeader(404)
			dev.ReportUserError(w, "Unknown RelationType")
			return
		}
	}

	for _, k := range ticket.Relations {
		//Cant be parent and child
		if (relation.RelationType == models.ParentOf && k.OtherTicket.ID == relation.OtherTicket && k.Type == models.ChildOf) || (relation.RelationType == models.ChildOf && k.OtherTicket.ID == relation.OtherTicket && k.Type == models.ParentOf) {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Can't be parent and child")
			return
		}

		if relation.OtherTicket == k.OtherTicket.ID && relation.RelationType == k.Type {
			w.WriteHeader(406)
			dev.ReportUserError(w, "Relation already exists")
			return
		}
	}

	id, err := db.AddRelation(ticketid, relation.OtherTicket, relation.RelationType)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(struct {
		Relation int64
	}{
		id,
	})
}

//DeleteRelation deletes a relation to a ticket
func DeleteRelation(w http.ResponseWriter, r *http.Request) {
	projectid, _ := strconv.ParseInt(strings.Split(r.RequestURI, "/")[4], 10, 64)
	queueid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[6], 10, 64)
	ticketid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[8], 10, 64)
	relationid, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[10], 10, 64)

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

	if !allperms.Admin && !perms.CanEditTicket {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to patch tickets in that queue.")
		return
	}

	ticket, found, err := db.GetTicket(ticketid, queueid, true)
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

	found = false

	for _, k := range ticket.Relations {
		if k.ID == relationid {
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Relation not found")
		return
	}

	err = db.DeleteRelation(relationid)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(struct {
		Relation string
	}{
		fmt.Sprintf("Relation %d deleted", relationid),
	})
}
