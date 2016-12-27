// Secure provides methods which provide secure
// handling of account and authentication data.
// Examples include hashing passwords, hmac, etc.
package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"github.com/cfrank/auth.fun/api/apierror"
	"golang.org/x/crypto/scrypt"
)

const (
	HASH_SALT_BYTES int = 32
	HASH_DK_BYTES   int = 64
)

// HashOptions are options provided to scrypt
type HashOptions struct {
	N       int // CPU/Mem cost param
	R       int // Block size param
	P       int // Parallelisation param
	SaltLen int // Length of salt
	DKLen   int // Length of derived key
}

// Some secure default hash options for scrypt
// http://www.tarsnap.com/scrypt/scrypt.pdf
var DefaultHashOptions = &HashOptions{
	N:       16384,
	R:       10,
	P:       10,
	SaltLen: HASH_SALT_BYTES,
	DKLen:   HASH_DK_BYTES,
}

// GenerateHashFromPassword provides a secure way of hashing a password with
// a sufficiently long salt for storage in a database. This function returns
// an array of:
// [N:R:P:Salt][Dkey]
func GenerateHashFromPassword(password string, options *HashOptions) ([]byte, *apierror.ApiError) {
	salt, saltError := GenerateSecureRandomBytes(options.SaltLen)

	if saltError != nil {
		return nil, apierror.ServerAuthError
	}

	scryptHash, scryptError := scrypt.Key([]byte(password), salt, options.N, options.R, options.P, options.DKLen)

	if scryptError != nil {
		return nil, apierror.ServerAuthError
	}

	// Hash the scrypt derived key with SHA256
	shaHash := sha256.New()
	shaHash.Write(scryptHash)

	return []byte(fmt.Sprintf("%s:%s:%d:%d:%d", encodeBytes(shaHash.Sum(nil)), encodeBytes(salt), options.N, options.R, options.P)), nil
}

// CompareHashToPassword compares in constant time a hash (generated from
// GenerateHashFromPassword) to a password provided from the frontend
// A nil return means the two byte slices have equal content
func CompareHashToPassword(hash, password []byte) *apierror.ApiError {
	dkey, salt, options, decodeError := decodeHash(hash)

	if decodeError != nil {
		return decodeError
	}

	compareDk, keyError := scrypt.Key(password, salt, options.N, options.R, options.P, options.DKLen)

	if keyError != nil {
		return apierror.ServerAuthError
	}

	if subtle.ConstantTimeCompare(hash, compareDk) == 1 {
		return nil
	}

	fmt.Println(dkey)

	return apierror.MismatchedHashError
}

// GenerateEmailVerifier generates a secure random string for use verifying a
// users email address.
func GenerateEmailVerifier() ([]byte, *apierror.ApiError) {
	verifier, verifierError := GenerateSecureRandomBytes(16)

	if verifierError != nil {
		return nil, apierror.ServerAuthError
	}

	return verifier, nil
}

// GenerateSecureRandomBytes generates a cryptographically secure
// random byte slice for use in secure functions
func GenerateSecureRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, randReadError := rand.Read(bytes)

	if randReadError != nil {
		return nil, randReadError
	}

	return bytes, nil
}

// decodeHash is a helper function which splits up the hash value
// generated from GenerateHashFromPassword and returns the derived key,
// the salt, and a pointer to HashOptions.
// The hash is stored together with the options so that theoretically
// every user could have different Scrypt params
func decodeHash(input []byte) ([]byte, []byte, *HashOptions, *apierror.ApiError) {
	decodedInput, decodeError := decodeString(string(input))

	if decodeError != nil {
		return nil, nil, nil, apierror.InvalidHashError
	}

	// [Dkey][Salt][N][R][P]
	var hashValues []string = strings.Split(string(decodedInput), ":")

	if len(hashValues) != 5 {
		return nil, nil, nil, apierror.InvalidHashError
	}

	var dkey, salt []byte
	var options *HashOptions = new(HashOptions)
	var err error

	dkey, err = decodeString(hashValues[0])

	if err != nil {
		return nil, nil, nil, apierror.InvalidHashError
	}

	salt, err = decodeString(hashValues[1])

	if err != nil {
		return nil, nil, nil, apierror.InvalidHashError
	}

	options.N, err = strconv.Atoi(hashValues[2])

	if err != nil {
		return nil, nil, nil, apierror.InvalidHashError
	}

	options.R, err = strconv.Atoi(hashValues[3])

	if err != nil {
		return nil, nil, nil, apierror.InvalidHashError
	}

	options.P, err = strconv.Atoi(hashValues[4])

	if err != nil {
		return nil, nil, nil, apierror.InvalidHashError
	}

	// Add length of dkey/salt to HashOptions
	options.DKLen = len(dkey)
	options.SaltLen = len(salt)

	return dkey, salt, options, nil
}

// encodeBytes is a simple helper function which encodes a byte slice
// into a base64 encoded string
func encodeBytes(input []byte) string {
	return base64.URLEncoding.EncodeToString(input)
}

// decodeString is a simple helper function which decodes a base64 encoded
// string and returns a byte slice or error
func decodeString(input string) ([]byte, error) {
	return base64.URLEncoding.DecodeString(input)
}
