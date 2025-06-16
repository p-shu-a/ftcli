package send

import (
	"context"
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

func SendFile(ctx context.Context, wg *sync.WaitGroup, file *os.File, rip net.IP, password string) {
	defer wg.Done()
	defer file.Close()

	dstConn, err := dialRemote(rip)
	if err != nil {
		log.Fatalf("failed to connect to remote ip: %v", err) /// log fatal or just a return?
	}
	defer dstConn.Close()

	hash, err := shared.FileChecksumSHA265(file)
	if err != nil {
		log.Fatal(err) /// log fatal or just a return?
	}

	salt, err := encryption.GenerateSalt()
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := encryption.GenerateNonce()
	if err != nil {
		log.Fatal(err)
	}

	header := models.Header{
		FileName: file.Name(),
		CheckSum: hash,
		Nonce:    nonce,
		Salt:     salt,
	}

	hdrJsonBytes, err := shared.HeaderToJsonB(header)
	if err != nil {
		log.Fatal(err)
	}

	hdrLen := shared.GetHeaderLength(hdrJsonBytes)

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	// create a byte-slice of the same size as the file to contain the plaintext
	plaintext := make([]byte, fileInfo.Size())
	if _, err := io.ReadFull(file, plaintext); err != nil {
		log.Fatal(err)
	}

	// Generate the cipher text
	cipherText, err := encryption.EncryptAEAD(salt, nonce, password, plaintext, hdrJsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	// send the header
	if err := sendHeader(dstConn, hdrJsonBytes, hdrLen); err != nil {
		log.Printf("failed to send header: %v", err)
	}

	// send the cipher text
	bytesWritten, err := dstConn.Write(cipherText)
	if err != nil {
		log.Fatalf("failed to send file: %v", err)
	}

	log.Printf("wrote %d bytes to remote", bytesWritten)
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
