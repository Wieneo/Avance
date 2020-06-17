package worker

import (
	//Imported to be used with database/sql
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
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
	smtpAlive := true
	err := db.Connection.QueryRow(`SELECT "Name" FROM "Patches" LIMIT 1`).Scan(&dummyDBVersion)
	if err != nil {
		dBAlive = false
		errors = append(errors, err.Error())
	}
	if !dBAlive {
		w.WriteHeader(500)
	}

	auth := smtp.PlainAuth("", config.CurrentConfig.SMTP.User, config.CurrentConfig.SMTP.Password, config.CurrentConfig.SMTP.Host)
	// TLS config
	tlsconfig := &tls.Config{
		ServerName: config.CurrentConfig.SMTP.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", config.CurrentConfig.SMTP.Host, config.CurrentConfig.SMTP.Port), tlsconfig)
	if err != nil {
		dev.LogError(err, fmt.Sprintf("Couldn't connect to SMTP server: %s", err.Error()))
		smtpAlive = false
		w.WriteHeader(500)
	} else {
		c, err := smtp.NewClient(conn, config.CurrentConfig.SMTP.Host)
		if err != nil {
			dev.LogError(err, fmt.Sprintf("Couldn't connect to SMTP server: %s", err.Error()))
			smtpAlive = false
			w.WriteHeader(500)
		} else {
			c.StartTLS(tlsconfig)

			// Auth
			if err = c.Auth(auth); err != nil {
				dev.LogError(err, fmt.Sprintf("Couldn't authenticate to SMTP server: %s", err.Error()))
				smtpAlive = false
				w.WriteHeader(500)
			} else {
				c.Close()
			}
		}
	}

	json.NewEncoder(w).Encode(struct {
		DB     bool
		SMTP   bool
		Errors []string
	}{
		dBAlive,
		smtpAlive,
		errors,
	})
}
