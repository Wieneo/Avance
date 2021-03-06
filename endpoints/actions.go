package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/templates"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

type actionWebRequest struct {
	Type    models.ActionType
	Content string
}

//CreateAction creates an action for a ticket
func CreateAction(w http.ResponseWriter, r *http.Request) {
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
		dev.ReportUserError(w, "You are not allowed to edit tickets in that queue.")
		return
	}

	_, found, err = db.GetTicket(ticketid, queueid, models.WantedProperties{})
	if err != nil {
		utils.ReportInternalErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, templates.TicketNotFound)
		return
	}

	var action actionWebRequest

	rawBytes, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(rawBytes, &action)
	if err != nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "JSON body is malformed")
		return
	}

	//Only comments and answers can be created via the API
	if action.Type != models.Comment && action.Type != models.Answer {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Only answers and comments can be created via the API")
		return
	}

	if utils.IsEmpty(action.Content) {
		w.WriteHeader(406)
		dev.ReportUserError(w, "Content can't be empty")
		return
	}

	title := "Comment was added"
	if action.Type == models.Answer {
		title = "Answer was added"
	}

	id, err := db.AddAction(ticketid, action.Type, title, action.Content, models.Issuer{Valid: true, Issuer: user})
	if err != nil {
		utils.ReportInternalErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Action int64
	}{
		id,
	})
}
