package auth

import (
	"encoding/base64"
	"encoding/json"
	//"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/cfrank/auth.fun/api/account"
	"github.com/cfrank/auth.fun/api/apierror"
)

type AuthData struct {
	DataString  string   `json:"authString"` // Base64 encoded auth data
	DataContent []string // Array of [username,password]
	Date        int64    `json:"date"`  // Unix time in seconds
	Nonce       string   `json:"nonce"` // Psuedo-Random 12 byte string
}

func NewAuth(w http.ResponseWriter, req *http.Request, params map[string]string) {
	// Make sure auth body was sent
	if req.Body == nil {
		err := apierror.New("The body of the request was empty", http.StatusBadRequest)
		err.Handle(w)
		return
	}

	defer req.Body.Close()

	// Create an empty AuthData struct to be populated
	var data *AuthData = new(AuthData)

	deconstructError := data.deconstructAuthString(req.Body)

	if deconstructError != nil {
		deconstructError.Handle(w)
		return
	}

	accountData, createAccountError := data.createAccount()

	if createAccountError != nil {
		createAccountError.Handle(w)
		return
	}

	emailVerifier, verifierError := GenerateEmailVerifier()

	if verifierError != nil {
		verifierError.Handle(w)
		return
	}

	accountData.EmailVerifier = emailVerifier

	// Database additions
	saveError := accountData.Save()

	if saveError != nil {
		saveError.Handle(w)
		return
	}
}

// Deconstruct the auth string from the client
// It comes in the format:
// {
//      authString: base64(username:password)
//      date: Unix time in seconds
//      nonce: Psuedo-Random 12 byte string
// }
func (data *AuthData) deconstructAuthString(input io.ReadCloser) *apierror.ApiError {
	jsonError := json.NewDecoder(input).Decode(&data)

	if jsonError != nil {
		return apierror.New("Malformed JSON received", http.StatusBadRequest)
	}

	// Deconstruct authString
	decodedAuthString, decodeError := base64.StdEncoding.DecodeString(data.DataString)

	if decodeError != nil {
		return apierror.New("Malformed authString received", http.StatusBadRequest)
	}

	data.DataContent = strings.Split(string(decodedAuthString), "&")

	if len(data.DataContent) != 2 {
		return apierror.New("Invalid authString received", http.StatusBadRequest)
	}

	return nil
}

// CreateAccount creates the users account in the database
// When the user is first created they are set up with a false
// flag for the verifiedEmail db table so they will need to
// do that before that is set
func (data *AuthData) createAccount() (*account.Account, *apierror.ApiError) {
	// Check for unique email address
	emailData, emailError := checkEmail(data.DataContent[0])

	if emailError != nil {
		return nil, emailError
	}

	accountId, accountIdError := account.GenerateAccountId()

	if accountIdError != nil {
		return nil, accountIdError
	}

	passwordHash, hashError := GenerateHashFromPassword(data.DataContent[1], DefaultHashOptions)

	if hashError != nil {
		return nil, hashError
	}

	var accountData *account.Account = new(account.Account)

	accountData.UserId = accountId
	accountData.EmailLocal = emailData[0]
	accountData.EmailDomain = emailData[1]
	accountData.PasswordHash = encodeBytes(passwordHash)
	accountData.EmailVerified = false

	return accountData, nil
}

// checkEmail takes an email address and splits it into two parts
// Example: chris@cfrank.org
// Local: chris
// Domain: cfrank.org
// Validates that these two parts of the correct length
// and returns them.
func checkEmail(email string) ([]string, *apierror.ApiError) {
	parts := strings.Split(email, "@")

	if len(parts) != 2 || len(email) > 254 || len(parts[0]) > 64 || len(parts[1]) > 190 || !strings.Contains(parts[1], ".") {
		return nil, apierror.New("Invalid email received", http.StatusBadRequest)
	}

	return parts, nil
}
