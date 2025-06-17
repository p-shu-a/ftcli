package send

import (
	"context"
	"fmt"
	"ftcli/config"
	"ftcli/internal/encryption"
	"ftcli/internal/shared"
	"ftcli/models"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func SendFile(ctx context.Context, wg *sync.WaitGroup, file *os.File, rip net.IP, password string) error {
	defer wg.Done()
	defer file.Close()

	dstConn, err := dialRemote(rip)
	if err != nil {
		return fmt.Errorf("failed to connect to remote ip: %v", err)
	}
	defer dstConn.Close()

	config.Dlog.Print("send: pt1")
	shared.PrintMemUsage()

	hash, err := shared.FileChecksumSHA265(file)
	if err != nil {
		return fmt.Errorf("failed to generate sha256 checksum of file: %v", err)
	}

	salt, err := encryption.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}
	nonce, err := encryption.GenerateNonce()
	if err != nil {
		return fmt.Errorf("failed to generate nonce: %v", err)
	}

	header := models.Header{
		FileName: file.Name(),
		CheckSum: hash,
		Nonce:    nonce,
		Salt:     salt,
	}

	hdrJsonBytes, err := shared.HeaderToJsonB(header)
	if err != nil {
		return fmt.Errorf("failed to convert header to json: %v", err)
	}

	hdrLen := shared.GetHeaderLength(hdrJsonBytes)

	config.Dlog.Print("send: pt2")
	shared.PrintMemUsage()

	fileInfo, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed to get file info: %v", err)
	}

	// create a byte-slice of the same size as the file to contain the plaintext
	plaintext := make([]byte, fileInfo.Size(), fileInfo.Size()+16)
	if _, err := io.ReadFull(file, plaintext); err != nil {
		return fmt.Errorf("failed to read from file: %v", err)
	}

	config.Dlog.Printf("plaintext size: %v\n", len(plaintext))
	config.Dlog.Printf("plaintext cap: %v\n", cap(plaintext))

	config.Dlog.Print("send: pt3. after readiing plaintext")
	shared.PrintMemUsage()

	// Generate the cipher texts
	cipherText, err := encryption.EncryptAEAD(salt, nonce, password, plaintext, hdrJsonBytes)
	if err != nil {
		return fmt.Errorf("failed to encrypt file: %v", err)
	}

	config.Dlog.Printf("ciphertext size: %v\n", len(cipherText))
	config.Dlog.Print("send: pt4. after generrating ciphertext")
	shared.PrintMemUsage()

	// send the header
	if err := sendHeader(dstConn, hdrJsonBytes, hdrLen); err != nil {
		return fmt.Errorf("failed to send header: %v", err)
	}

	config.Dlog.Print("send: pt5")
	shared.PrintMemUsage()

	// send the cipher text
	bytesWritten, err := dstConn.Write(cipherText)
	if err != nil {
		return fmt.Errorf("failed to write encypted file contents to remote: %v", err)
	}

	config.Dlog.Print("send: pt6.end")
	shared.PrintMemUsage()
	log.Printf("wrote %d bytes to remote", bytesWritten)

	return nil
}

// Dials the remote address and returns a connection (net.Conn)
func dialRemote(rip net.IP) (net.Conn, error) {

	remoteAddr := net.TCPAddr{
		IP:   rip,
		Port: config.ReceivePort,
	}
	conn, err := net.Dial("tcp", remoteAddr.String())
	if err != nil {
		return nil, err
	}
	return conn, nil

}

// This function writes the header to the peer's conn
// Takes the net.conn to the peer, header in the form of json-encoded bytes, and the length of the header in bytes.
func sendHeader(dst net.Conn, hdrJsonBytes []byte, hdrLen []byte) error {

	// Let peer know how large the header they're about to receive is
	if _, err := dst.Write(hdrLen); err != nil {
		return err
	}
	// send actual header
	if _, err := dst.Write(hdrJsonBytes); err != nil {
		return err
	}
	return nil

}
