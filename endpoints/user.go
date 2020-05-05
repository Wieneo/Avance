package endpoints

import (
	"encoding/json"
	"net/http"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/models"

	"gitlab.gnaucke.dev/tixter/tixter-app/v2/dev"
)

//GetProfile returns the profile of the currently logged in user to the client
func GetProfile(w http.ResponseWriter, r *http.Request) {
	user, err := GetUser(r, w)
	if err != nil {
		w.WriteHeader(500)
		dev.ReportError(w, err.Error())
	}

	json.NewEncoder(w).Encode(struct {
		models.User
	}{
		user,
	})
}