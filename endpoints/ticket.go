package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"
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
		json.NewEncoder(w).Encode(struct {
			Ticket models.Ticket
		}{
			ticket,
		})
	} else {
		w.WriteHeader(401)
		dev.ReportUserError(w, "You don't have access to that queue!")
	}
}
