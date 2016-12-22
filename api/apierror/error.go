// Provides a simple way to throw errors that are encountered
// when trying to access api endpoints
// When an error is handled it responds to the client with the
// correct error code and a message to display
package apierror

import (
	"encoding/json"
	"net/http"
)

type ApiError struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

var (
	ServerAuthError *ApiError = New("The server encountered an error processing the request", http.StatusInternalServerError)

	SessionIdError *ApiError = New("Invalid Session ID received", http.StatusBadRequest)

	MismatchedHashError *ApiError = New("Mismatched Hash received", http.StatusUnauthorized)

	InvalidHashError *ApiError = New("An invalid Hash value was recieved", http.StatusBadRequest)
)

// New returns a initialized ApiError
func New(message string, errorCode int) *ApiError {
	return &ApiError{
		Message:   message,
		ErrorCode: errorCode,
	}
}

// Handle handles the error and sends the information back to the client
func (err *ApiError) Handle(w http.ResponseWriter) {
	w.WriteHeader(err.ErrorCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
