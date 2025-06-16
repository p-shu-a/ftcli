package encryption

import (
	"crypto/cipher"
	"net"

	"golang.org/x/crypto/chacha20"
)

// This file includes ChaCha20 related encryption/decryption functions

// Returns a new chacha20 cipher stream
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
