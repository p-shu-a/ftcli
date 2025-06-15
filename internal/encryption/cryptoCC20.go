package encryption

import (
	"crypto/cipher"
	"crypto/rand"
	"ftcli/config"
	"net"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20"
)


// Returns a 16byte salt
func GenerateSalt() ([]byte, error) {

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil

}

// Returns a 12byte nonce
func GenerateNonce() ([]byte, error) {

	nonce := make([]byte, chacha20.NonceSize)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	return nonce, nil

}

// Generates a master key
func GenerateMasterKey(salt []byte, password string) []byte {
	return argon2.Key(
		[]byte(password),
		salt,
		config.Time,
		config.Memory,
		config.Threads,
		config.KeyLength)
}

// Returns a new chacha20 cipher stream. Key is generated here.
func newChaCha20CipherStream(salt []byte, nonce []byte, password string) (*chacha20.Cipher, error) {

	masterKey := GenerateMasterKey(salt, password)
	stream, err := chacha20.NewUnauthenticatedCipher(masterKey, nonce)
	if err != nil {
		return nil, err
	}
	return stream, err

}

// Return a StreamWriter for encrypting with ChaCha20
func EncryptSetupChaCha20(salt []byte, nonce []byte, password string, dstConn net.Conn) (*cipher.StreamWriter, error) {

	stream, err := newChaCha20CipherStream(salt, nonce, password)
	if err != nil {
		return nil, err
	}
	strWriter := &cipher.StreamWriter{
		S: stream,
		W: dstConn,
	}
	return strWriter, nil

}

// Return a StreamReader for Decrypting with ChaCha20
func DecryptSetupChaCha20(salt []byte, nonce []byte, password string, srcConn net.Conn) (*cipher.StreamReader, error) {

	stream, err := newChaCha20CipherStream(salt, nonce, password)
	if err != nil {
		return nil, err
	}
	strReader := &cipher.StreamReader{
		S: stream,
		R: srcConn,
	}
	return strReader, nil

}