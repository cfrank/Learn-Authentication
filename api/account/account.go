package account

import (
	"math/rand"

	"github.com/cfrank/auth.fun/api/apierror"
	"github.com/cfrank/auth.fun/api/database"
)

type Account struct {
	UserId        string // 16 bytes
	EmailLocal    string
	EmailDomain   string
	PasswordHash  []byte
	EmailVerified bool
}

const (
	ACCOUNT_ID_LEN   int    = 16
	LETTER_BYTES     string = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ-_"
	MAX_ACC_ID_TRIES int    = 5
)

func GenerateAccountId() (string, *apierror.ApiError) {
	for i := 0; i < MAX_ACC_ID_TRIES; i++ {
		accountId := string(psuedoRandomBytes(ACCOUNT_ID_LEN))

		if database.UniqueAccountId(accountId) == true {
			return accountId, nil
		}
	}

	return "", apierror.ServerAuthError
}

func psuedoRandomBytes(length int) []byte {
	bytes := make([]byte, length)
	for i := range bytes {
		bytes[i] = LETTER_BYTES[rand.Int63()%int64(len(LETTER_BYTES))]
	}

	return bytes
}
