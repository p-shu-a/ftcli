package send

import (
	"context"
	"encoding/binary"
	"errors"
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

	config.Slog.Print("send: pt1")
	shared.PrintMemUsage()

	hash, err := shared.FileChecksumSHA265(file)
	if err != nil {
		return fmt.Errorf("failed to generate sha256 checksum of file: %v", err)
	}

	// Sending InfoHeader
	// Filename, Checksum, and Salt are sent via info-header. They are only sent once.
	salt, err := encryption.GenerateSalt()
	if err != nil {
		return fmt.Errorf("failed to generate salt: %v", err)
	}
	nonce, err := encryption.GenerateNonce()
	if err != nil {
		return fmt.Errorf("failed to generate nonce: %v", err)
	}
	baseNonce := nonce[:4]
	masterKey := encryption.GenerateMasterKey(salt, password)
	info := models.Header{
		FileName: file.Name(),
		CheckSum: hash,
		Salt:     salt,
		Nonce:    baseNonce,
	}
	infoHdrBytes, _ := shared.HeaderToJsonB(info)
	infoLen := shared.GetHeaderLength(infoHdrBytes)

	if err := sendHeader(dstConn, infoHdrBytes, infoLen); err != nil {
		return err
	}

	// Total bytes written to dst
	var totalBytesWritten int = 0
	// Since the encrypted files are sent in chunks, keep track of how many chunks are sent
	var chunkIndex int = 0
	// plaintext is read from file in chunks that are FileChunkSize long
	plaintextChunk := make([]byte, config.FileChunkSize)

	// This loop repeatedly chunks the plaintext data, sends a related header, and then sends the encypted data
	for {

		// Read from file. If file can't be read, proceed no further
		////// what if there is stuff still in the buffer from before? how do you clear that?
		n, err := file.Read(plaintextChunk)

		// following ifs need improvement
		if err != nil {
			if errors.Is(err, io.EOF) { // n would be 0 here
				break
			}
			return err
		}

		/// unique nonce every time, or same once  once + chunkIdx?
		binary.BigEndian.PutUint64(nonce[4:], uint64(chunkIndex))
		if err != nil {
			return fmt.Errorf("failed to generate nonce: %v", err)
		}

		// Package the nonce into the header
		header := models.Header{
			Nonce: nonce,
		}
		hdrJsonBytes, err := shared.HeaderToJsonB(header)
		if err != nil {
			return fmt.Errorf("failed to convert header to json: %v", err)
		}
		hdrLen := shared.GetHeaderLength(hdrJsonBytes)

		// Generate the cipherText
		cipherText, err := encryption.EncryptAEAD(nonce, masterKey, plaintextChunk[:n], hdrJsonBytes)
		if err != nil {
			return fmt.Errorf("failed to encrypt file: %v", err)
		}

		// send the header
		if err := sendHeader(dstConn, hdrJsonBytes, hdrLen); err != nil {
			return fmt.Errorf("failed to send header: %v", err)
		}

		// send the cipher text
		bytesWritten, err := dstConn.Write(cipherText)
		if err != nil {
			return fmt.Errorf("failed to write encypted file contents to remote: %v", err)
		}
		totalBytesWritten += bytesWritten

		chunkIndex++
	}
	log.Printf("wrote %d bytes to remote", totalBytesWritten)
	config.Slog.Print("send: pt6.end")
	shared.PrintMemUsage()

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

// This function writes the header to the peer's conn.
// It sends: the length of the header and then the header contents
// Params are: net.conn to the peer, length of header in bytes, and header in the form of json-encoded bytes
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
