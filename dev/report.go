package dev

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
)

//ReportError sends back a error message to the user
func ReportError(Error error, w http.ResponseWriter, Message string) {
	json.NewEncoder(w).Encode(struct {
		Error string
	}{
		Message,
	})

	sentry.CaptureException(Error)
	sentry.Flush(time.Second * 5)
}

//ReportUserError sends back a error message to the user without informing sentry
func ReportUserError(w http.ResponseWriter, Message string) {
	json.NewEncoder(w).Encode(struct {
		Error string
	}{
		Message,
	})

}
