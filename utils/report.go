package utils

import (
	"net/http"

	"gitlab.gnaucke.dev/avance/avance-app/v2/dev"
)

//ReportErrorToUser is used to send an internal server error message to the user
func ReportErrorToUser(Error error, w http.ResponseWriter) {
	w.WriteHeader(500)
	dev.ReportError(Error, w, "An error happened on our side :( Please try again later!")
	//Sentry already gets informed via dev.ReportError
	dev.LogErrorNoSentry(Error, Error.Error())
}
