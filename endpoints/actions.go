package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
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

	_, found, err = db.GetTicket(ticketid, true)
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

	id, err := db.AddAction(ticketid, action.Type, title, action.Content, user.ID)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	json.NewEncoder(w).Encode(struct {
		Action int64
	}{
		id,
	})
}
