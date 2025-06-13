package receive

import (
	"context"
	"encoding/binary"
	"encoding/json"
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

// receive mode is relitively simple. just opens a listener and waits to get
func ReceiveFile(ctx context.Context, wg *sync.WaitGroup, password string) {
	defer wg.Done()

	localLn := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: config.ReceivePort,
	}

	ln, err := net.ListenTCP("tcp", &localLn)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	} else {
		log.Printf("Waiting for file transfer on port: %d", localLn.Port)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// what are the possilbe errors?
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		log.Printf("Receiving on port: %d", localLn.Port)
		go downloadFile(conn, password)
	}
}

// Download the actual file
func downloadFile(srcConn net.Conn, password string) {

	log.Printf("Downloading file from %v....", srcConn.RemoteAddr().String())

	// Receive the header
	hdr, err := receiveHeader(srcConn)
	if err != nil {
		log.Printf("failed to receive header: %v", err)
		return
	}
	log.Printf("Filename: %v", hdr.FileName)
	log.Printf("sha256 checksum: %v", hdr.CheckSum)

	// Allow user to accept or decline
	log.Printf("Continue Download? (yes/no/y/n/Y/N): ")
	var resp string
	fmt.Scanln(&resp)
	log.Printf("%q", resp)
	if resp != "yes" {
		srcConn.Close()
		log.Print("Not downloading...")
		return
	}

	strReader, err := encryption.DecryptSetup(hdr.IV, password, srcConn)
	if err != nil {
		log.Printf("Decryption setup failed: %v", err)
		return
	}

	// create the file for saving
	file, err := os.Create(hdr.FileName)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return
	}
	defer file.Close()

	// write to file and get hash of file
	hash, bytesReceived, err := shared.CopyAndHash(file, strReader)
	if err != nil {
		log.Printf("error copying file: %v", err)
	}
	log.Printf("successfully copied %d bytes", bytesReceived)
	log.Printf("hash of file is : %v", hash)

	if hash == hdr.CheckSum {
		log.Printf("hashes match")
	} else {
		log.Printf("Hash mismatch!")
	}
	log.Printf("finished downloading %v", hdr.FileName)

}

// Receive the header
func receiveHeader(conn net.Conn) (models.Header, error) {

	var hdr models.Header
	var lenBuf [4]byte

	// read header length
	if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
		return models.Header{}, err
	}

	// turn header length into a unsigned (hum0n readable) int
	headerLen := binary.BigEndian.Uint32(lenBuf[:])

	log.Printf("STATS FOR NERDS: header len: %v", headerLen)

	// read the jsonbytes
	jsonBytes := make([]byte, headerLen)
	if _, err := io.ReadFull(conn, jsonBytes); err != nil {
		return models.Header{}, err
	}

	// turn into a struct
	if err := json.Unmarshal(jsonBytes, &hdr); err != nil {
		return models.Header{}, err
	}

	return hdr, nil
}
