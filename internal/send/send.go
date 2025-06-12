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

	conn := dialRemote(rip)

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
func dialRemote(rip net.IP) net.Conn {
	remoteAddr := net.TCPAddr{
		IP: rip,
		Port: config.ReceivePort,
	}
	conn, err := net.Dial("tcp", remoteAddr.String())
	if err != nil{
		log.Fatalf("failed to dial remote: %v", err)
	}
	return conn
}

func sendHeader(conn net.Conn, header models.Header) error {

	jsonData, err := json.Marshal(header)
	if err != nil{
		return err
	}
	var lenBuf [4]byte
	binary.BigEndian.PutUint32(lenBuf[:], uint32(len(jsonData)))

	if _, err := conn.Write(lenBuf[:]); err != nil{
		return err
	}

	if _ , err := conn.Write(jsonData); err != nil {
		return err
	}

	return nil
}