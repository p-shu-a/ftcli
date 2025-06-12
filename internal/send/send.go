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

func SendFile(ctx context.Context, wg *sync.WaitGroup, fileInfo os.FileInfo, rip net.IP){
	defer wg.Done()
	
	log.Printf("file name: %v", fileInfo.Name())
	// calculate sha256 checksum of file
	f, err := os.Open(fileInfo.Name())
	if err != nil{
		log.Fatalf("failed to read file: %v", err)
	}
	defer f.Close()

	hash, err := shared.FileChecksum(f)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("file hash is: %v", hash)

	conn := dialRemote(rip)

	header := models.Header{
		FileName: f.Name(),
		CheckSum: hash,
	}
	if err := sendHeader(conn, header); err != nil{
		log.Printf("failed to send header: %v", err)
	}
	
	bytesWritten, err := io.Copy(conn, f)
	if err != nil{
		log.Fatalf("failed to send file %v", err)
	}

	log.Printf("wrote %d bytes to remote", bytesWritten)
	
}

// Dials the remote address and returns a connection (net.Conn)
func dialRemote(rip net.IP) net.Conn {
	remoteAddr := net.TCPAddr{
		IP: rip,
		Port: config.RecievePort,
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