// VerifiyEmail provides methods for dealing with verifiying a users
// email address.
package auth

type EmailVerification struct {
	AccountId string
	Verifier  []byte
	Expires   int64
}

func (data *EmailVerification) VerifyEmail() {
	// Save verification details in DB
	// Send Email verification
}
