package encryption

import (
	"golang.org/x/crypto/chacha20poly1305"
)

// This file includes AEAD related encryption/decryption functions

// Encrypt the plaintext using ChaCha20 and generate a Poly1305 MAC. The Mac is sealed with the cipher text.
func EncryptAEAD(salt []byte, nonce []byte, password string, plaintext []byte, adHeader []byte) ([]byte, error) {

	key := GenerateMasterKey(salt, password)

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}

	cipherText := aead.Seal(plaintext[:0], nonce, plaintext, adHeader)

	return cipherText, err

}

// Authenticates the ciphertext and, if legit, decrypts it, and returns the plaintext
// The decrypted plaintext is stored in the passed cipherText slice (mostly) replacing it
func DecryptAEAD(salt []byte, nonce []byte, password string, cipherText []byte, jsonHdrBytes []byte) ([]byte, error) {

	key := GenerateMasterKey(salt, password)

	aead, err := chacha20poly1305.New(key)
	if err != nil {
		return nil, err
	}
	plaintext, err := aead.Open(cipherText[:0], nonce, cipherText, jsonHdrBytes)
	if err != nil {
		return nil, err
	}

	return plaintext, nil

}
