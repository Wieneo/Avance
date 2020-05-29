package worker

import (
	//Imported to be used with database/sql
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
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

func getInstanceHealth(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(struct {
		Status string
	}{
		"Worker is alive!",
	})
}
