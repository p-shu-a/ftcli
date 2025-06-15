package send

import (
	"context"
	"encoding/binary"
	"encoding/json"
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

	hash, err := shared.FileChecksumSHA265(file)
	if err != nil {
		log.Fatal(err) /// log fatal or just a return?
	}

	dstConn, err := dialRemote(rip)
	if err != nil {
		log.Fatalf("failed to connect to remote ip: %v", err) /// log fatal or just a return?
	}


	//////////////////////////////////// AES business /////////////////////////////////////
	// iv, err := encryption.GenerateIV()
	// if err != nil {
	// 	log.Fatal(err) /// log fatal or just a return?
	// }

	// // Create a header
	// header := models.HeaderAES{
	// 	FileName: file.Name(),
	// 	CheckSum: hash,
	// 	IV:       iv,
	// }

	// // get a cipher stream
	// strWriter, err := encryption.EncryptSetupAES(header.IV, password, dstConn)
	// if err != nil {
	// 	log.Fatal(err) /// log fatal or just a return?
	// }
	//////////////////////////////////// AES business END /////////////////////////////////////

	//////////////////////////////////// ChaCha20 business /////////////////////////////////////
	salt, err := encryption.GenerateSalt()
	if err != nil{
		log.Fatal(err)
	}
	nonce, err := encryption.GenerateNonce()
	if err != nil {
		log.Fatal(err)
	}
	strWriter, err := encryption.EncryptSetupChaCha20(salt, nonce, password, dstConn)
	if err != nil {
		log.Fatal(err)
	}

	header := models.Header{
		FileName: file.Name(),
		CheckSum: hash,
		Nonce: nonce,
		Salt: salt,
	}
	//////////////////////////////////// AES business /////////////////////////////////////

	if err := sendHeader(dstConn, header); err != nil {
		log.Printf("failed to send header: %v", err)
	}

	// encrypt and send data
	bytesWritten, err := io.Copy(strWriter, file)
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

// Sends header to peer
func sendHeader(dst net.Conn, header models.Header) error {

	// marshal to json formatted bytes
	headerBytes, err := json.Marshal(header)
	if err != nil {
		return err
	}
	// calcualte the length of the header
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(headerBytes)))

	// Let peer know how large the header they're about to receive is
	if _, err := dst.Write(lenBuf[:]); err != nil {
		return err
	}
	// send actual header
	if _, err := dst.Write(headerBytes); err != nil {
		return err
	}

	return nil
}
