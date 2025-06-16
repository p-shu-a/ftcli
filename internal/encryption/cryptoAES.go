package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"net"
)

// This file includes AES encryption/decryption related functions

/*
This function returns the cipher stream. On the way there it:
creates a password hash as the key,
a new cipher block,
converts the block to a stream using newCTR.
Encrypt and Decrypt both share this functionality
*/
func newCipherStreamAES(iv []byte, password string) (cipher.Stream, error) {

	// Get sha256 sum of password
	key := sha256.Sum256([]byte(password))
	// generate new cipher block
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}
	// get a stream cipher from this block cipher! oh i get it!
	stream := cipher.NewCTR(block, iv)
	return stream, nil

}

// Helper function to encrypt data for sending
func EncryptSetupAES(iv []byte, password string, dstConn net.Conn) (*cipher.StreamWriter, error) {

	cipherStr, err := newCipherStreamAES(iv, password)
	if err != nil {
		return nil, err
	}
	writer := &cipher.StreamWriter{S: cipherStr, W: dstConn}
	return writer, nil

}

// Helper Function to decrypt data upon receiving
func DecryptSetupAES(iv []byte, password string, srcConn net.Conn) (*cipher.StreamReader, error) {

	cipherStr, err := newCipherStreamAES(iv, password)
	if err != nil {
		return nil, err
	}
	reader := &cipher.StreamReader{S: cipherStr, R: srcConn}
	return reader, nil

}
