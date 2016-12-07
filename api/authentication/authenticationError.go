package authentication

import (
	"encoding/json"
	"net/http"
)

type AuthenticationError struct {
	Message   string `json:"message"`
	ErrorCode int    `json:"code"`
}

var (
	ServerAuthError *AuthenticationError = &AuthenticationError{
		Message:   "The server encountered an error processing the request",
		ErrorCode: http.StatusInternalServerError,
	}

	SessionIdError *AuthenticationError = &AuthenticationError{
		Message:   "Invalid Session ID received",
		ErrorCode: http.StatusBadRequest,
	}

	MismatchedHashError *AuthenticationError = &AuthenticationError{
		Message:   "Mismatched Hash received",
		ErrorCode: http.StatusUnauthorized,
	}
)

func (err *AuthenticationError) AuthenticationError(w http.ResponseWriter) {
	w.WriteHeader(err.ErrorCode)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(err)
}
