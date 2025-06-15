package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"net"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20"
)

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

// Generates the initilization vector (IV)
// Returns IV (16bytes)
func GenerateIV() ([]byte, error) {

	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}
	return iv, nil

}

////////////////// ChaCha20 Operations //////////////////


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

func DecryptSetupChaCha20(salt []byte, nonce []byte, password string, srcConn net.Conn) (*cipher.StreamReader, error){

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

// Returns a new chacha20 cipher stream
func newChaCha20CipherStream(salt []byte, nonce []byte, password string) (*chacha20.Cipher, error) {
	key := argon2.Key([]byte(password), salt, 3, 32*1024, 4, 32)
	
	stream, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil{
		return nil, err
	}
	
	return stream, err
}


// Returns a 16byte salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _ , err := rand.Read(salt); err != nil{
		return nil, err
	}
	return salt, nil
}

// Returns a 12byte nonce
func GenerateNonce() ([]byte, error) {
	nonce := make([]byte, chacha20.NonceSize)
	if _, err := rand.Read(nonce); err != nil{
		return nil, err
	}
	return nonce, nil
}
















// // shit method name
// func EncryptSetupCC20P1305(password string) (error) {
// 	/*	Things I need
// 		- key generated from password (32bytes)
// 		- nonce (12bytes)
// 		- generate a MAC, in another program

// 	*/
// 	salt := make([]byte, 16)
// 	key := argon2.Key([]byte(password), salt, 3, 32*1024, 4, 32)

// 	_, err := chacha20poly1305.New(key)
// 	if err != nil {
// 		return err
// 	}
	
// 	nonce := make([]byte, chacha20poly1305.NonceSize)
// 	if _, err := rand.Read(nonce); err != nil{
// 		return err
// 	}
// 	return nil
// }