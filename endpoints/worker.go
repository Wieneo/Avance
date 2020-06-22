package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/perms"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//GetWorkers returns all registered workers
func GetWorkers(w http.ResponseWriter, r *http.Request) {
	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	userperms, err := perms.CombinePermissions(user)

	if !userperms.Admin && !userperms.CanSeeWorker {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to view all workers")
		return
	}

	if workers, err := db.GetAllWorkers(); err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
	} else {
		json.NewEncoder(w).Encode(workers)
	}
}

//ToggleWorker enables or disables a worker
func ToggleWorker(w http.ResponseWriter, r *http.Request) {
	workerID, _ := strconv.Atoi(strings.Split(r.URL.String(), "/")[4])

	user, err := utils.GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	userperms, err := perms.CombinePermissions(user)

	if !userperms.Admin && !userperms.CanChangeWorker {
		w.WriteHeader(403)
		dev.ReportUserError(w, "You are not allowed to enalbe/disable workers")
		return
	}

	worker, found, err := db.GetWorker(workerID)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
		return
	}

	if !found {
		w.WriteHeader(404)
		dev.ReportUserError(w, "Specified worker wasn't found")
		return
	}

	worker.Active = !worker.Active
	if err := db.PatchWorker(worker); err != nil {
		w.WriteHeader(500)
		dev.ReportError(err, w, err.Error())
	} else {
		json.NewEncoder(w).Encode(worker)
	}
}
