// Package cryptography is the package which hosts encrypting and decrypting functions
//
// Exported function(s): EncryptMessage, DecryptMessage
package cryptography

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"log"
)

// EncryptMessage encrypts a string using AES
//
// Parameters:
//
//	key ([]byte) - The key to use for encryption
//	message (string) - The message to encrypt
//
// Returns: (string,error) where string is the encoded string.
// On success, error should be nil
func EncryptMessage(key []byte, message string) (string, error) {
	msgToByte := []byte(message)
	// Generate a cipher using AES, the length of the input key determines the version of AES used:
	// Key length of 16 bytes -> AES-128
	// Key length of 24 bytes -> AES-192
	// Key length of 32 bytes -> AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Could not create new cipher: %v", err)
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(msgToByte))
	iv := cipherText[:aes.BlockSize]
	// Check that cipher text is of correct length
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		log.Printf("Could not encrypt message: %v", err)
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], msgToByte)

	// Encode the cipher text in base64
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

// DecryptMessage decrypts a string using AES
//
// Parameters:
// key ([]byte) - The key to use for decryption
// message (string) - The message to decrypt
//
// Returns: (string,error) where string is the decrypted message.
// On success, error should be nil
func DecryptMessage(key []byte, message string) (string, error) {
	cipherText, err := base64.StdEncoding.DecodeString(message)
	if err != nil {
		log.Printf("Could not decode base64 message: %v", err)
		return "", err
	}
	// For information on AES key length, see the EncryptMessage function
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Printf("Could not create new cipher: %v", err)
		return "", err
	}

	if len(cipherText) < aes.BlockSize {
		log.Printf("Invalid cipher text length: %v", err)
		return "", err
	}
	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
