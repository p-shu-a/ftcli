package receive

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
	"strings"
	"sync"
)

// Creates a listener and gets user input on whether to accept file download or not
// If user accepts download, calls downloadFile
func ReceiveFile(ctx context.Context, wg *sync.WaitGroup, password string) error {
	defer wg.Done()

	localLn := net.TCPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: config.ReceivePort,
	}

	ln, err := net.ListenTCP("tcp", &localLn)
	if err != nil {
		return fmt.Errorf("failed to create listener: %v", err)
	}

	for {
		log.Printf("Receiving on port: %d", localLn.Port)
		conn, err := ln.Accept()
		if err != nil {
			// what are the possible errors?
			if errors.Is(err, net.ErrClosed) {
				return err
			}
			continue
		}
		inforHdr, err := receiveHeader(conn) // perhaps receiveHeader ad JsonBToHeader should be consolidated
		if err != nil {
			return err
		}
		hdr, err := shared.JsonBToHeader(inforHdr)
		if err != nil {
			return err
		}

		log.Printf("Download from  : %v", conn.RemoteAddr().String())
		log.Printf("Filename       : %v", hdr.FileName)
		log.Printf("SHA256 checksum: %v", hdr.CheckSum)

		// Allow user to accept or decline
		log.Printf("Continue Download? (yes/no): ")
		var resp string
		fmt.Scanln(&resp)
		resp = strings.ToLower(resp)
		if resp != "yes" {
			// conn.Close() /// should i close the conn here?
			log.Print("Not downloading...")
			continue
		} else {
			// not too sure it needs to be a Go func
			go downloadFile(conn, password, hdr)
		}
	}
}

// Download the actual file
func downloadFile(srcConn net.Conn, password string, infoHdr *models.Header) {

	config.Dlog.Printf("Downloading file from %v....", srcConn.RemoteAddr().String())

	// Do some file name validation
	var file *os.File
	var ctr int = 1
	for {
		if _, err := os.Stat(infoHdr.FileName); err != nil {
			// File does not exists, create it
			file, err = os.Create(infoHdr.FileName)
			if err != nil {
				log.Printf("failed to create file: %v", err)
				return
			}
			defer file.Close()
			break
		} else {
			// Original file already exists.
			// Find a new name for the file
			newName := shared.SuggestNewFileName(infoHdr.FileName, ctr)

			// Check if new name already exists
			if _, err := os.Stat(newName); err != nil {
				log.Printf("File with name %v already exists. Setting new name: %v", infoHdr.FileName, newName)
				file, err = os.Create(newName)
				if err != nil {
					log.Printf("failed to create file: %v", err)
					return
				}
				defer file.Close()
				break
			} else {
				// new selected name also exists, try again loop
				ctr++
				continue
			}
		}
	}

	masterkey := encryption.GenerateMasterKey(infoHdr.Salt, password)
	cipherText := make([]byte, config.FileChunkSize+16)
	// keeps track of chunks received
	var chunkCtr int = 0

	for {

		// receive header from sending peer
		hdrJsonBytes, err := receiveHeader(srcConn)
		if err != nil {
			if errors.Is(err, io.EOF) {
				config.Dlog.Printf("got EOF from remote, breaking")
				break
			}
			log.Printf("failed to receive header: %v", err)
			return
		}

		hdr, err := shared.JsonBToHeader(hdrJsonBytes)
		if err != nil {
			log.Fatal(err)
		}

		// From remote read cipherText
		n, err := srcConn.Read(cipherText)
		if err != nil {
			log.Fatal(err)
		}

		// Decrypt the cipherText
		plaintext, err := encryption.DecryptAEAD(hdr.Nonce, masterkey, cipherText[:n], hdrJsonBytes)
		if err != nil {
			log.Printf("failed to drcrypt: %v", err)
			// if there is an issue decrypting, close and remove the created file
			file.Close()
			os.Remove(file.Name())
			return
		}

		// Write plaintext to file
		_, err = file.Write(plaintext)
		if err != nil {
			log.Printf("failed to write plaintext to file: %v", err)
			return
		}

		chunkCtr++
	}

	config.Dlog.Print("receive : end")
	shared.PrintMemUsage()

	log.Printf("finished downloading %v", file.Name())
	hash, err := shared.FileChecksumSHA265(file)
	if err != nil {
		log.Printf("failed to generate hash for file: %v", err)
		return
	}

	log.Printf("hash from infoHdr: %v", infoHdr.CheckSum)
	log.Printf("hash of DLed File: %v", hash)
	if hash == infoHdr.CheckSum {
		log.Printf("hashes match")
	} else {
		log.Printf("Hash mismatch!")
	}
}

// Helper to read header coming from remote, clean it up and return header json bytes
// First comes the header length. Then comes the header. Read header-len from conn to get header
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
