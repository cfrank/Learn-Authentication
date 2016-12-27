package account

import (
	//"fmt"
	"math/rand"

	"github.com/cfrank/auth.fun/api/apierror"
	"github.com/cfrank/auth.fun/api/database"
)

type Account struct {
	UserId        string // 16 bytes
	EmailLocal    string // Max 64 bytes
	EmailDomain   string // Max 190 bytes
	PasswordHash  string // Max 140 bytes (if scrypt settings change)
	EmailVerifier []byte // 16 bytes
	EmailVerified bool
}

const (
	ACCOUNT_ID_LEN   int    = 16
	LETTER_BYTES     string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	MAX_ACC_ID_TRIES int    = 5
)

func (account *Account) Save() *apierror.ApiError {
	// Save Account into database
	insertError := database.InsertAccount(account.UserId, account.EmailLocal, account.EmailDomain, account.PasswordHash, account.EmailVerified)

	if insertError != nil {
		return insertError
	}

	return nil
}

func GenerateAccountId() (string, *apierror.ApiError) {
	for i := 0; i < MAX_ACC_ID_TRIES; i++ {
		accountId := string(psuedoRandomBytes(ACCOUNT_ID_LEN))

		if database.UniqueAccountId(accountId) == true {
			return accountId, nil
		}
	}

	// Could not find an available account id
	return "", apierror.ServerAuthError
}

func psuedoRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = LETTER_BYTES[rand.Int63()%int64(len(LETTER_BYTES))]
	}

	return bytes
}
