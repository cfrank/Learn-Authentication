package authentication

const (
	// Max difference (in seconds) in time between request
	// and retrieval by the server
	MAXTIMEDIFF int = 30

	// Length of scrypt salt in bytes
	SC_SALT_BYTES int = 32

	// Length of scrypt derived key in bytes
	SC_DK_BYTES int = 64
)
