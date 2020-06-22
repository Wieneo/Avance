package worker

import (
	//Imported to be used with database/sql
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.gnaucke.dev/avance/avance-app/v2/config"
	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
)

//StartServing starts serving the health endpoint
func StartServing() {
	dev.LogInfo("Listening on 0.0.0.0:", strconv.Itoa(config.CurrentConfig.Port))
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	router.HandleFunc("/api/v1/health", getInstanceHealth).Methods("GET")

	srv := &http.Server{
		Handler: router,
		Addr:    "0.0.0.0:" + strconv.Itoa(config.CurrentConfig.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	srv.ListenAndServe()
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		dev.LogDebug(r.Method, r.RemoteAddr, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//getInstanceHealth returns the current state of the instance
//This is mainly used to check if requests function
//Output should look something like this: {"DB":true,"Redis":true,"Errors":[]}
func getInstanceHealth(w http.ResponseWriter, r *http.Request) {
	errors := make([]string, 0)

	var dummyDBVersion string

	dBAlive := true
	err := db.Connection.QueryRow(`SELECT "Name" FROM "Patches" LIMIT 1`).Scan(&dummyDBVersion)
	if err != nil {
		dBAlive = false
		errors = append(errors, err.Error())
	}
	if !dBAlive {
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(struct {
		DB     bool
		Errors []string
	}{
		dBAlive,
		errors,
	})
}
