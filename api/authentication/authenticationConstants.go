package authentication

const (
	// Max difference (in seconds) in time between request
	// and retrieval by the server
	MAXTIMEDIFF int = 30

	// Selector length
	TOKEN_SELECTOR_BYTES int = 12

	// Length of scrypt salt
	SC_SALT_BYTES int = 32

	// Length of scrypt derived key
	SC_DK_BYTES int = 64
)
