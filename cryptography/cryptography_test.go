package cryptography

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"

// TestEncryptDecrypt tests the EncryptMessage and DecryptMessage from the cryptography package
//
// Parameters:
//
//	t (*testing.T) - A reference to the testing framework
//
// Returns: _
func TestEncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")
	message := "Hello"

	encrypted, err := EncryptMessage(key, message)
	require.Nil(t, err)
	require.Regexp(t, isBase64, encrypted)

	decrypted, err := DecryptMessage(key, encrypted)
	require.Nil(t, err)
	require.Equal(t, message, decrypted)

	// Test encryption with an invalid key,
	// key is invalid because it is not 16, 24, or 32 bytes long
	key = []byte("0123456789abcde")
	_, err = EncryptMessage(key, message)
	require.NotNil(t, err)

	// Test decryption with an invalid message,
	// message is invalid because it is not a base64 string
	_, err = DecryptMessage(key, "invalid")
	require.NotNil(t, err)
	// Test decryption with an invalid key,
	_, err = DecryptMessage(key, encrypted)
}
