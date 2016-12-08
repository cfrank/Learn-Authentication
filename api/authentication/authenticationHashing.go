package authentication

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/scrypt"
)

type ScryptParams struct {
	N       int // CPU/Mem cost param
	R       int // Block size param
	P       int // Parallelisation param
	SaltLen int // Bytes to use in salt
	DKLen   int // Length of derived key
}

var DefaultScryptParams = &ScryptParams{
	N:       16384,
	R:       8,
	P:       1,
	SaltLen: SC_SALT_BYTES,
	DKLen:   SC_DK_BYTES,
}

/*
 * Generates a securely random byte slice
 *
 * Takes a integer len which is the length of the random
 * byte slice
 */
func GenerateRandomBytes(len int) ([]byte, error) {
	bytes := make([]byte, len)
	_, randReadError := rand.Read(bytes)

	if randReadError != nil {
		return nil, randReadError
	}

	return bytes, nil
}

/*
 * Generate a auth token
 *
 * The auth token is contains the following information
 * Selector - A random string used to query the auth_token
 * table for a hashed authentication validator
 */
func GenerateAuthSelector(len int) ([]byte, *AuthenticationError) {
	authSelector, selectorError := GenerateRandomBytes(len)

	if selectorError != nil {
		return nil, ServerAuthError
	}

	return []byte(encodeBytes(authSelector)), nil
}

/*
 * Generate a hash from a byte slice
 *
 * Returns a byte slice in the format
 * [MEM COST]&[BlOCK SIZE]&[PAR PARAM]&[SALT]&[DK]
 * where the salt and dk are base16 lowercase (2 bytes per byte)
 */
func GenerateHashFromSlice(password []byte, params *ScryptParams) ([]byte, *AuthenticationError) {
	salt, saltError := GenerateRandomBytes(params.SaltLen)

	if saltError != nil {
		return nil, ServerAuthError
	}

	dkey, keyError := scrypt.Key(password, salt, params.N, params.R, params.P, params.DKLen)

	if keyError != nil {
		return nil, ServerAuthError
	}

	return []byte(fmt.Sprintf("%d&%d&%d&%s&%s", params.N, params.R, params.P, encodeBytes(salt), encodeBytes(dkey))), nil
}

/*
 * Compare the hash with a provided password
 */
func CompareHashToPassword(hash []byte, password []byte) *AuthenticationError {
	salt, dk, params, decodeError := DecodeHashString(hash)

	if decodeError != nil {
		return decodeError
	}

	compareDk, keyError := scrypt.Key(password, salt, params.N, params.R, params.P, params.DKLen)

	if keyError != nil {
		return ServerAuthError
	}

	if subtle.ConstantTimeCompare(dk, compareDk) == 1 {
		return nil
	}

	return MismatchedHashError
}

/*
 * Decode the hash and return the values
 */
func DecodeHashString(hash []byte) ([]byte, []byte, ScryptParams, *AuthenticationError) {
	// First split the string
	var hashValues []string = strings.Split(string(hash), "&")

	// Make sure the correct number of values are present
	if len(hashValues) != 5 {
		return nil, nil, ScryptParams{}, SessionIdError
	}

	// Decode the Scrypt parameters
	var params ScryptParams
	var err error

	params.N, err = strconv.Atoi(hashValues[0])

	if err != nil {
		return nil, nil, params, SessionIdError
	}

	params.R, err = strconv.Atoi(hashValues[1])

	if err != nil {
		return nil, nil, params, SessionIdError
	}

	params.P, err = strconv.Atoi(hashValues[2])

	if err != nil {
		return nil, nil, params, SessionIdError
	}

	salt, err := base64.URLEncoding.DecodeString(hashValues[3])

	if err != nil {
		return nil, nil, params, SessionIdError
	}
	params.SaltLen = len(salt)

	dk, err := base64.URLEncoding.DecodeString(hashValues[4])

	if err != nil {
		return nil, nil, params, SessionIdError
	}
	params.DKLen = len(dk)

	return salt, dk, params, nil
}

/*
 * A helper function which returns a base64 encoded string
 * for use when we need a url safe string
 */
func encodeBytes(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}
