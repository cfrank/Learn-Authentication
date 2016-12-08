package authentication

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

/*
 * Struct which holds the incoming authentication request
 *
 * AuthenticationData - Holds the username and password base64 encoded
 * Date - Holds the seconds since Epoch
 * Nonce - Psuedo random string provided by front-end which will be returned
 * so front-end can validate response
 */
type AuthenticationData struct {
	AuthenticationData string `json:"authString"`
	Date               int64  `json:"date"`
	Nonce              string `json:"nonce"`
}

/*
 * Creates a new authentication profile for a users data
 */
func NewAuth(w http.ResponseWriter, req *http.Request, params map[string]string) {
	var authData *AuthenticationData = &AuthenticationData{}

	// Get information from the request
	if req.Body == nil {
		err := &AuthenticationError{
			Message:   "Invalid request - No request body received",
			ErrorCode: http.StatusBadRequest,
		}
		err.AuthenticationError(w)
		return
	}

	jsonError := json.NewDecoder(req.Body).Decode(&authData)

	if jsonError != nil {
		err := &AuthenticationError{
			Message:   "Malformed JSON recieved",
			ErrorCode: http.StatusBadRequest,
		}
		err.AuthenticationError(w)
		return
	}

	reqError := authData.validateRequest()

	if reqError != nil {
		reqError.AuthenticationError(w)
		return
	}

	authValues, reqStringError := authData.deconstructAuthData()

	if reqStringError != nil {
		reqStringError.AuthenticationError(w)
		return
	}

	// Generate a authentication token selector
	authSelector, selectorError := GenerateAuthSelector(TOKEN_SELECTOR_BYTES)

	if selectorError != nil {
		selectorError.AuthenticationError(w)
		return
	}

	// Generate hash value for password
	hash, hashError := GenerateHashFromSlice([]byte(authValues[1]), DefaultScryptParams)

	if hashError != nil {
		ServerAuthError.AuthenticationError(w)
		return
	}

	fmt.Printf("%s:%s", string(authSelector), string(hash))
}

func (req *AuthenticationData) validateRequest() *AuthenticationError {
	// Make sure all fields are correctly populated
	if req.AuthenticationData == "" || req.Nonce == "" {
		err := &AuthenticationError{
			Message:   "Recieved empty authentication data",
			ErrorCode: http.StatusBadRequest,
		}
		return err
	}

	// Make sure the request was recieved in due time
	timeDiff := int(time.Now().Unix() - req.Date)
	if timeDiff > MAXTIMEDIFF || timeDiff < 0 {
		err := &AuthenticationError{
			Message:   "Invalid request date recieved, or the request timed out",
			ErrorCode: http.StatusBadRequest,
		}
		return err
	}

	return nil
}

func (req *AuthenticationData) deconstructAuthData() ([]string, *AuthenticationError) {
	// First decode base64
	decodedString, decodeError := base64.StdEncoding.DecodeString(req.AuthenticationData)
	if decodeError != nil {
		return nil, ServerAuthError
	}

	// decodedString should be in format email&password
	var authValues []string = strings.Split(string(decodedString), "&")

	// Make sure correct number of values recieved
	if len(authValues) != 2 {
		err := &AuthenticationError{
			Message:   "Invalid AuthenticationData recieved",
			ErrorCode: http.StatusBadRequest,
		}
		return nil, err
	}

	// Make sure the values are not empty
	if authValues[0] == "" || authValues[1] == "" {
		err := &AuthenticationError{
			Message:   "Empty Authentication values recieved",
			ErrorCode: http.StatusBadRequest,
		}
		return nil, err
	}

	return authValues, nil
}
