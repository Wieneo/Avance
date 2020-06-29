package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//GetTaskInfo returns info about a task
func GetTaskInfo(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	taskID, _ := strconv.ParseInt(strings.Split(r.URL.String(), "/")[4], 10, 64)

	task, found, err := db.GetTask(taskID)

	if err != nil {

		utils.ReportErrorToUser(err, w)
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Specified task wasn't found")
		return
	}

	if task.Ticket.Valid {
		ticket, _, err := db.GetTicketUnsafe(task.Ticket.Int64, models.WantedProperties{Queue: true})
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		allperms, perms, err := perms.GetPermissionsToQueue(user, ticket.Queue)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		if !perms.CanSee && !allperms.Admin {
			w.WriteHeader(401)
			dev.ReportUserError(w, "You are not authorized to view tasks in that queue")
			return
		}
	} else {
		allperms, err := perms.CombinePermissions(user)
		if err != nil {

			utils.ReportErrorToUser(err, w)
			return
		}

		if !allperms.Admin {
			w.WriteHeader(401)
			dev.ReportUserError(w, "You are not allowed to view global tasks")
			return
		}
	}

	json.NewEncoder(w).Encode(task)

}
