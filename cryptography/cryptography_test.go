// Testing file for the cryptography package
package cryptography

import (
	"github.com/stretchr/testify/require"
	"testing"
)

const isBase64 = "^(?:[A-Za-z0-9+/]{4})*(?:[A-Za-z0-9+/]{2}==|[A-Za-z0-9+/]{3}=|[A-Za-z0-9+/]{4})$"

// TestEncryptDecrypt tests the EncryptMessage and DecryptMessage from the cryptography package
func TestEncryptDecrypt(t *testing.T) {
	key := []byte("0123456789abcdef")
	message := "Hello"

	encrypted, err := EncryptMessage(key, message)
	require.Nil(t, err)
	require.Regexp(t, isBase64, encrypted)

	decrypted, err := DecryptMessage(key, encrypted)
	require.Nil(t, err)
	require.Equal(t, message, decrypted)
}
