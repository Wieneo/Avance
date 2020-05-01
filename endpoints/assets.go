package endpoints

import (
	"net/http"
)

//ServeAssets serves common assets to clients
func ServeAssets(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/app/dist/"+r.RequestURI)
}

//ServeAppFrontend serves the react app for the frontend users
func ServeAppFrontend(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./frontend/app/dist/index.html")
}
