package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/redis"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
)

//GetInstanceHealth returns the current state of the instance
//This is mainly used to check if requests function
//Output should look something like this: {"Version":"0.1 Dev","DB":true,"Redis":true,"Errors":[]}
func GetInstanceHealth(w http.ResponseWriter, r *http.Request) {
	errors := make([]string, 0)

	var dummyDBVersion string

	dBAlive := true
	err := db.Connection.QueryRow(`SELECT "Name" FROM "Patches" LIMIT 1`).Scan(&dummyDBVersion)
	if err != nil {
		dBAlive = false
		errors = append(errors, err.Error())
	}

	redisAlive, err := redis.Ping()
	if err != nil {
		errors = append(errors, err.Error())
	}

	if !dBAlive || !redisAlive {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(struct {
		Version string
		DB      bool
		Redis   bool
		Errors  []string
	}{
		fmt.Sprint(utils.AppVersion, " ", utils.AppChannel),
		dBAlive,
		redisAlive,
		errors,
	})
}
