package receive

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"ftcli/config"
	"ftcli/internal/encryption"
	"ftcli/internal/shared"
	"io"
	"log"
	"net"
	"os"
	"strings"
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
			// what are the possible errors?
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		log.Printf("Receiving on port: %d", localLn.Port)
		// Download accept/decline should be here
		go downloadFile(conn, password)
	}
}

// Download the actual file
func downloadFile(srcConn net.Conn, password string) {

	log.Printf("Downloading file from %v....", srcConn.RemoteAddr().String())

	// Receive the header
	hdrJsonBytes, err := receiveHeader(srcConn)
	if err != nil {
		log.Printf("failed to receive header: %v", err)
		return
	}

	hdr, err := shared.JsonBToHeader(hdrJsonBytes)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("Filename: %v", hdr.FileName)
	// log.Printf("sha256 checksum: %v", hdr.CheckSum)
	// log.Printf("entire header: %v", hdr)

	// Allow user to accept or decline
	log.Printf("Continue Download? (yes/no): ")
	var resp string
	fmt.Scanln(&resp)
	resp = strings.ToLower(resp)
	if resp != "yes" {
		srcConn.Close()
		log.Print("Not downloading...")
		return
	}

	// Since our sender only sends the header len, header and ciphertext and then closes the connection
	// we can simply read the rest of the bytes from the conn
	cipherText, err := io.ReadAll(srcConn)
	if err != nil {
		log.Fatal(err)
	}

	/// decrypt with AEAD
	plaintext, err := encryption.DecryptAEAD(hdr.Salt, hdr.Nonce, password, cipherText, hdrJsonBytes)
	if err != nil {
		log.Printf("Failed to drcrypt: %v", err)
		return
	}

	// create the file for saving
	file, err := os.Create(hdr.FileName) // what if file with same name already exists?
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return
	}
	defer file.Close()

	// write to file and get hash of file
	hash, bytesReceived, err := shared.CopyAndHash(file, bytes.NewReader(plaintext))
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
func receiveHeader(conn net.Conn) ([]byte, error) {

	// read header length
	var lenBuf [4]byte
	if _, err := io.ReadFull(conn, lenBuf[:]); err != nil {
		return nil, err
	}

	// turn header length into a unsigned int
	headerLen := binary.BigEndian.Uint32(lenBuf[:])

	// read the jsonbytes
	hdrJsonBytes := make([]byte, headerLen)
	if _, err := io.ReadFull(conn, hdrJsonBytes); err != nil {
		return nil, err
	}

	return hdrJsonBytes, nil

}
