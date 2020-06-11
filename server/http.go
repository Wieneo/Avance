package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
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
	router.HandleFunc("/api/v1/ping", servePong).Methods("GET")

	router.HandleFunc("/api/v1/login", endpoints.LoginUser).Methods("POST")

	router.HandleFunc("/api/v1/logout", endpoints.LogoutUser).Methods("GET")

	router.HandleFunc("/api/v1/profile", endpoints.GetProfile).Methods("GET")
	router.HandleFunc("/api/v1/profile/avatar", endpoints.UpdateProfilePicture).Methods("POST")
	router.HandleFunc("/api/v1/profile/avatar", endpoints.RemoveProfilePicture).Methods("DELETE")
	router.HandleFunc("/api/v1/profile/avatar", endpoints.GetProfilePicture).Methods("GET")
	router.HandleFunc("/api/v1/profile/settings", endpoints.PatchSettings).Methods("PATCH")
	router.HandleFunc("/api/v1/profile", endpoints.PatchProfile).Methods("PATCH")

	router.HandleFunc("/api/v1/users", endpoints.GetUsers).Methods("GET")
	router.HandleFunc("/api/v1/users", endpoints.CreateUser).Methods("POST")
	router.HandleFunc("/api/v1/groups", endpoints.GetGroups).Methods("GET")
	router.HandleFunc("/api/v1/groups", endpoints.CreateGroup).Methods("POST")

	router.HandleFunc("/api/v1/user/{[0-9]{*}}", endpoints.DeactivateUser).Methods("DELETE")
	router.HandleFunc("/api/v1/user/{[0-9]{*}}/avatar", endpoints.GetProfilePicture).Methods("GET")
	router.HandleFunc("/api/v1/user/{[0-9]{*}}/permissions", endpoints.GetPermissionsOfUser).Methods("GET")

	//PROJECT APIs
	router.HandleFunc("/api/v1/projects", endpoints.GetProjects).Methods("GET")
	router.HandleFunc("/api/v1/projects", endpoints.CreateProject).Methods("POST")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}", endpoints.GetSingleProject).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}", endpoints.ChangeProject).Methods("PATCH")

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

	//Queues
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queues", endpoints.GetProjectQueues).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queues", endpoints.CreateQueue).Methods("POST")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}", endpoints.PatchQueue).Methods("PATCH")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}", endpoints.DeleteQueue).Methods("DELETE")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/tickets", endpoints.GetTicketsFromQueue).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/tickets", endpoints.CreateTicketsInQueue).Methods("POST")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}", endpoints.GetTicketFullPath).Methods("GET")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}", endpoints.PatchTicketsInQueue).Methods("PATCH")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}/owner", endpoints.DeletePropertyFromTicket).Methods("DELETE")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}/stalleduntil", endpoints.DeletePropertyFromTicket).Methods("DELETE")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}/relations", endpoints.CreateRelation).Methods("POST")
	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}/relation/{[0-9]{*}}", endpoints.DeleteRelation).Methods("DELETE")

	router.HandleFunc("/api/v1/project/{[0-9]{*}}/queue/{[0-9]{*}}/ticket/{[0-9]{*}}/actions", endpoints.CreateAction).Methods("POST")

	router.HandleFunc("/api/v1/ticket/{[0-9]{*}}", endpoints.GetTicket).Methods("GET")

	router.HandleFunc("/api/v1/workers", endpoints.GetWorkers).Methods("GET")
	router.HandleFunc("/api/v1/worker/{[0-9]{*}}", endpoints.ToggleWorker).Methods("PATCH")

	//Needs to be at the bottom!
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
		dev.LogDebug(r.Method, r.RemoteAddr, r.RequestURI)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

//sitesForUnauthorized contains all URLs which should be accessible without being logged in
var sitesForUnauthorized = []string{
	"^\\/$",
	"^\\/\\?.*$",
	"^\\/settings.*$",
	"^\\/login*",
	"^\\/js\\/*",
	"^\\/css\\/*",
	"^\\/api\\/v1\\/session$",
	"^\\/api\\/v1\\/health$",
	"^\\/api\\/v1\\/ping$",
	"^\\/api\\/v1\\/login$",
}

//authorizationMiddleware gets called at every request to check if user is authenticated
func authorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !canBeIgnored(r.RequestURI) {
			if !IsAuthorized(r) {
				if strings.HasPrefix(r.RequestURI, "/api/") {
					w.WriteHeader(401)
					json.NewEncoder(w).Encode(struct {
						Error string
					}{
						"You are currently not authorized!",
					})
				}
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

//servePong only returns pong. Used to check if instance is alive.
func servePong(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Pong!"))
}
