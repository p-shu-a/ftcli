package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"net"
)

// Helper function to encrypt data for sending
func EncryptSetup(iv []byte, password string, dstConn net.Conn) (*cipher.StreamWriter, error) {

	cipherStr, err := newCipherStream(iv, password)
	if err != nil {
		return nil, err
	}
	writer := &cipher.StreamWriter{S: cipherStr, W: dstConn}
	return writer, nil

}

// Helper Function to decrypt data upon receiving
func DecryptSetup(iv []byte, password string, srcConn net.Conn) (*cipher.StreamReader, error) {

	cipherStr, err := newCipherStream(iv, password)
	if err != nil {
		return nil, err
	}
	reader := &cipher.StreamReader{S: cipherStr, R: srcConn}
	return reader, nil

}

/*
This function returns the cipher stream. On the way there it:
creates a password hash,
a new cipher block,
converts the block to a stream using newCTR.
Encrypt and Decrypt both share this functionality
*/
func newCipherStream(iv []byte, password string) (cipher.Stream, error) {
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

// Generates the initilization vector (IV)
// Returns IV (16bytes)
func GenerateIV() ([]byte, error) {

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil

}
