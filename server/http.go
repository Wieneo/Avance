package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/perms"

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

	router.Use(perms.CheckAccessToProject)

	router.HandleFunc("/api/v1/health", endpoints.GetInstanceHealth).Methods("GET")
	router.HandleFunc("/api/v1/session", serveSessionInfo).Methods("GET")
	router.HandleFunc("/api/v1/login", endpoints.LoginUser).Methods("POST")

	router.HandleFunc("/api/v1/logout", endpoints.LogoutUser).Methods("GET")

	router.HandleFunc("/api/v1/profile", endpoints.GetProfile).Methods("GET")

	//PROJECT APIs
	router.HandleFunc("/api/v1/projects", endpoints.GetProjects).Methods("GET")
	router.HandleFunc("/api/v1/projects", endpoints.CreateProject).Methods("POST")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}", endpoints.GetSingleProject).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}", endpoints.ChangeProject).Methods("PATCH")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queues", endpoints.GetProjectQueues).Methods("GET")

	//Severities
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/severities", endpoints.GetSeverities).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/severities", endpoints.CreateSeverity).Methods("POST")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/severity/{[0-9]{*}}", endpoints.PatchSeverity).Methods("PATCH")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/severity/{[0-9]{*}}", endpoints.DeleteSeverity).Methods("DELETE")

	//Statuses
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/statuses", endpoints.GetStatuses).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/statuses", endpoints.CreateStatus).Methods("POST")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/status/{[0-9]{*}}", endpoints.PatchStatus).Methods("PATCH")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/status/{[0-9]{*}}", endpoints.DeleteStatus).Methods("DELETE")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/tickets", endpoints.GetTicketsFromQueue).Methods("GET")
	router.HandleFunc("/api/v1/ticket/{[0-9]{*}}", endpoints.GetTicket).Methods("GET")

	//Needs to be at the bottom!
	router.HandleFunc("/", endpoints.ServeAppFrontend).Methods("GET")
	router.HandleFunc("/login", endpoints.ServeAppFrontend).Methods("GET")
	router.PathPrefix("/").HandlerFunc(endpoints.ServeAssets).Methods("GET")

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
	"/api/v1/health",
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
		//Check if maybe cookie was set
		keks, err := r.Cookie("session")
		if err != nil {
			return false
		}

		session = keks.Value
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
