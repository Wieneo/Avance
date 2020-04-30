package endpoints

import (
	"encoding/json"
	"net/http"
)

//GetInstanceHealth returns the current state of the instance
//This is mainly used to check if requests function
func GetInstanceHealth(w http.ResponseWriter, r *http.Request) {
	//ToDo: Proper health page
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(struct {
		Status string
	}{
		"Instance is responding!",
	})
}
