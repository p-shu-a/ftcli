package send

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"ftcli/config"
	"ftcli/internal/shared"
	"ftcli/models"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

func SendFile(ctx context.Context, wg *sync.WaitGroup, file *os.File, rip net.IP){
	defer wg.Done()
	defer file.Close()

	log.Printf("file name: %v", file.Name())
	
	// calculate sha256 checksum of file
	hash, err := shared.FileChecksum(file)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("file hash: %v", hash)

	conn, err := dialRemote(rip)
	if err != nil{
		log.Fatalf("failed to connect to remote ip: %v", err)
	}

	header := models.Header{
		FileName: file.Name(),
		CheckSum: hash,
	}
	if err := sendHeader(conn, header); err != nil{
		log.Printf("failed to send header: %v", err)
	}
	
	bytesWritten, err := io.Copy(conn, file)
	if err != nil{
		log.Fatalf("failed to send file: %v", err)
	}

	log.Printf("wrote %d bytes to remote", bytesWritten)
	
}

// Dials the remote address and returns a connection (net.Conn)
func dialRemote(rip net.IP) (net.Conn, error) {
	remoteAddr := net.TCPAddr{
		IP: rip,
		Port: config.ReceivePort,
	}
	conn, err := net.Dial("tcp", remoteAddr.String())
	if err != nil{
		return nil, err
	}
	return conn, nil
}

func sendHeader(dst net.Conn, header models.Header) error {

	jsonData, err := json.Marshal(header)
	if err != nil{
		return err
	}
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(jsonData)))

	if _, err := dst.Write(lenBuf[:]); err != nil{
		return err
	}

	if _ , err := dst.Write(jsonData); err != nil {
		return err
	}

	return nil
}