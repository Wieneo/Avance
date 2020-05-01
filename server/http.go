package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/config"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/endpoints"
	"gitlab.gnaucke.dev/tixter/tixter-app/v2/redis"

	"github.com/gorilla/mux"
)

//HTTPInit starts the http endpoint
func HTTPInit() {
	dev.LogInfo("Listening on 0.0.0.0:", strconv.Itoa(config.CurrentConfig.Port))
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	router.Use(authorizationMiddleware)

	router.HandleFunc("/api/v1/health", endpoints.GetInstanceHealth)
	router.HandleFunc("/api/v1/session", serveSessionInfo)

	//Needs to be at the bottom!
	router.HandleFunc("/", endpoints.ServeAppFrontend)
	router.HandleFunc("/login", endpoints.ServeAppFrontend)
	router.PathPrefix("/").HandlerFunc(endpoints.ServeAssets)

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
		// Do stuff here
		dev.LogDebug(r.RemoteAddr, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//sitesForUnauthorized contains all URLs which should be accessible without being logged in
var sitesForUnauthorized = []string{
	"/login$",
	"/js/*",
	"/css/*",
	"/api/v1/session",
}

//authorizationMiddleware gets called at every request to check if user is authenticated
func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !canBeIgnored(r.RequestURI) {
			if !IsAuthorized(r) {
				/*w.WriteHeader(401)
				json.NewEncoder(w).Encode(struct {
					Error string
				}{
					"You are currently not authorized!\nPlease log-in first",
				})*/
				http.Redirect(w, r, "/login", 302)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func canBeIgnored(URL string) bool {
	for _, k := range sitesForUnauthorized {
		if result, err := regexp.MatchString(k, URL); err == nil && result {
			return true
		}
	}
	return false
}

//IsAuthorized returns true if user is authorized right now and refreshes the session token
func IsAuthorized(r *http.Request) bool {
	session := r.Header.Get("Authorization")
	if len(session) == 0 {
		return false
	}

	return redis.SessionValid(session)
}

//serveSessionInfo tells the client if the session key is working
func serveSessionInfo(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct {
		Authorized bool
	}{
		IsAuthorized(r),
	})
}
