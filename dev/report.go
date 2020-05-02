package dev

import (
	"encoding/json"
	"net/http"
)

//ReportError sends back a error message to the user
func ReportError(w http.ResponseWriter, Message string) {
	json.NewEncoder(w).Encode(struct {
		Error string
	}{
		Message,
	})
}
