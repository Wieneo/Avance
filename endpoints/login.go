package endpoints

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"gitlab.gnaucke.dev/avance/avance-app/v2/db"
	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
	"gitlab.gnaucke.dev/avance/avance-app/v2/redis"
	"gitlab.gnaucke.dev/avance/avance-app/v2/utils"
	"golang.org/x/crypto/bcrypt"
)

//LoginUser is called when a user send a POST Request to /api/v1/login
func LoginUser(w http.ResponseWriter, r *http.Request) {
	//Check if a body was sent with the request
	if r.Body == nil {
		w.WriteHeader(400)
		dev.ReportUserError(w, "You must send a body with your request!")
		return
	}

	var loginRequest struct {
		Username string
		Password string
	}

	//Try to parse body into LoginRequest struct
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &loginRequest); err != nil {
		w.WriteHeader(400)
		dev.ReportError(err, w, "Invalid Body recieved. Please check your JSON.")
		return
	}

	//Select ID + Password from Database
	if rows, err := db.Connection.Query(`SELECT "ID","Password" FROM "Users" WHERE "Username" = $1 AND "Active" = true`, loginRequest.Username); err != nil {
		utils.ReportErrorToUser(err, w)
	} else {
		//If the query returned an empty result set
		if !rows.Next() {
			w.WriteHeader(401)
			dev.ReportUserError(w, "Combination of username and password doesn't match")
			return
		}

		var UserID int64
		var PasswordHash string

		if err := rows.Scan(&UserID, &PasswordHash); err != nil {
			utils.ReportErrorToUser(err, w)
			return
		}

		rows.Close()

		//Compare password from request with hashed password from query
		if correct := bcrypt.CompareHashAndPassword([]byte(PasswordHash), []byte(loginRequest.Password)); correct != nil {
			w.WriteHeader(401)
			dev.ReportUserError(w, "Combination of username and password doesn't match")
			return
		}

		//Generate Session ID
		SessionKey, err := redis.CreateSession(UserID)
		if err != nil {
			utils.ReportErrorToUser(err, w)
			return
		}

		//Send back new session key
		json.NewEncoder(w).Encode(struct {
			SessionKey string
		}{
			SessionKey,
		})
	}

}

//LogoutUser is called when a user send a POST Request to /api/v1/login
func LogoutUser(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("Authorization")
	if len(session) == 0 {
		//Check if maybe cookie was set
		keks, err := r.Cookie("session")
		if err != nil {
			utils.ReportErrorToUser(err, w)
			return
		}

		session = keks.Value
	}

	if len(session) == 0 {
		w.WriteHeader(404)
		dev.ReportUserError(w, "No session found")
		return
	}

	if err := redis.DestroySession(r); err != nil {
		utils.ReportErrorToUser(err, w)
		return
	}

	json.NewEncoder(w).Encode(struct {
		Result string
	}{
		"Session destroyed",
	})
}
