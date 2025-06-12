package send

import (
	"context"
	"ftcli/config"
	"ftcli/internal/shared"
	"log"
	"net"
	"os"
	"sync"
)

func SendFile(ctx context.Context, wg *sync.WaitGroup, file os.FileInfo, rip net.IP){
	defer wg.Done()
	
	log.Printf("file name: %v", file.Name())
	// calculate sha256 checksum of file
	f, err := os.Open(file.Name())
	if err != nil{
		log.Fatalf("failed to read file: %v", err)
	}
	defer f.Close()
	
	// serialize
	conn := dialRemote(rip)
	// send
	checkSum, bytesWritten, err := shared.CopyAndHash(conn, f)
	if err != nil{
		log.Printf("error while copying file: %v", err)
	}
	log.Printf("file checksum is: %v", checkSum)
	log.Printf("wrote %d bytes to remote", bytesWritten)
	
}


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

