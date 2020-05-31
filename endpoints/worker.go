package endpoints

import (
	"encoding/json"
	"net/http"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/db"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/utils"
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
