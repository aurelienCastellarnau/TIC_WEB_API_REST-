package helpers

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// ErrorResponse API definition type (cf consigne)
type ErrorResponse struct {
	Code          int    `json:"code,omitempty"`
	Message       string `json:"message,omitempty"`
	Missingfields string `json:"missingFields,omitempty"`
}

// ErrorResponseJSON exported []byte ErrorResponse to write in response
var ErrorResponseJSON, _ = json.Marshal(ErrorResponse{
	Code:    http.StatusBadRequest,
	Message: "ErrorResponse",
})

// MustBeAdminJSON exported []byte ErrorResponse to write in response
var MustBeAdminJSON, _ = json.Marshal(ErrorResponse{
	Code:    http.StatusForbidden,
	Message: "Must be admin",
})

// MustBeConnectedJSON exported []byte ErrorResponse to write in response
var MustBeConnectedJSON, _ = json.Marshal(ErrorResponse{
	Code:    http.StatusUnauthorized,
	Message: "Must be connected",
})

// ParseErr custom error type
type ParseErr struct {
	Message string `json:"message,string"`
}

func (err ParseErr) Error() string {
	return err.Message
}

// AuthError status 401 send authentification error headers
func AuthError(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Basic realm=\"Must be connected\"")
	http.Error(w, "Must be connected", http.StatusUnauthorized)
}

// SendErrorResponse status 400 write errorResponse json in response and log message.
func SendErrorResponse(w http.ResponseWriter, message string, err error) {
	if err != nil && err.Error() != "" {
		fmt.Printf(message+"%s", err.Error())
	} else {
		fmt.Printf(message)
	}
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintf(w, "%s", ErrorResponseJSON)
}

// SendMustBeConnectedResponse status 401 write errorResponse json in response and log message.
func SendMustBeConnectedResponse(w http.ResponseWriter, message string, err error) {
	if err != nil && err.Error() != "" {
		fmt.Printf(message+"%s", err.Error())
	} else {
		fmt.Printf(message)
	}
	w.WriteHeader(http.StatusUnauthorized)
	fmt.Fprintf(w, "%s", MustBeConnectedJSON)
}

// SendMustBeAdminResponse status 403 write errorResponse json in response and log message.
func SendMustBeAdminResponse(w http.ResponseWriter, message string, err error) {
	if err != nil && err.Error() != "" {
		fmt.Printf(message+"%s", err.Error())
	} else {
		fmt.Printf(message)
	}
	w.WriteHeader(http.StatusForbidden)
	fmt.Fprintf(w, "%s", MustBeAdminJSON)
}

// SendNotFoundErrorResponse status 404 write errorResponse json in response and log message.
func SendNotFoundErrorResponse(w http.ResponseWriter, message string, err error) {
	if err != nil && err.Error() != "" {
		fmt.Printf(message+"%s", err.Error())
	} else {
		fmt.Printf(message)
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s", ErrorResponseJSON)
}

// SendUnauthorisedResponse status 401 write unauthorized header
func SendUnauthorisedResponse(w *http.ResponseWriter, message string, err error) {
	fmt.Printf(message+"%s", err.Error())
	http.Error(*w, "Must be connected", http.StatusUnauthorized)
}
