package endpoints

import (
	"net/http"
	"strings"
)

var frontendEndpoints = []string{
	"/?",
	"/login",
	"/settings",
}

//ServeAssets serves common assets to clients
func ServeAssets(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/" || matchesFrontendEndpoints(r.RequestURI) {
		http.ServeFile(w, r, "./frontend/app/dist/index.html")
	} else {
		http.ServeFile(w, r, "./frontend/app/dist/"+r.RequestURI)
	}
}

func matchesFrontendEndpoints(requestURI string) bool {
	for _, k := range frontendEndpoints {
		if strings.HasPrefix(requestURI, k) {
			return true
		}
	}
	return false
}
