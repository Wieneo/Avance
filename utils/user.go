package utils

import (
	"errors"
	"net/http"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/models"
	"gitlab.gnaucke.dev/avance/avance-app/v2/redis"
)

//GetUserID gets the user id from a web request
func GetUserID(r *http.Request) (int64, error) {
	session := r.Header.Get("Authorization")
	if len(session) == 0 {
		//Check if maybe cookie was set
		keks, err := r.Cookie("session")
		if err != nil {
			return 0, errors.New("No user assigned to request")
		}

		session = keks.Value
	}

	id, err := redis.SessionToUserID(session)
	if err != nil {
		return 0, err
	}
	return id, nil
}

//GetUser is a shortcut to get the assigned user to a request
func GetUser(r *http.Request, w http.ResponseWriter) (models.User, error) {
	userid, err := GetUserID(r)
	if err != nil {
		ReportInternalErrorToUser(err, w)
		return models.User{}, err
	}

	user, _, err := db.GetUser(userid)
	if err != nil {
		ReportInternalErrorToUser(err, w)
		return models.User{}, err
	}

	return user, nil

}
