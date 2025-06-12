package recieve

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"ftcli/config"
	"ftcli/internal/shared"
	"ftcli/models"
	"io"
	"log"
	"net"
	"os"
	"sync"
)

// Recieve mode is relitively simple. just opens a listener and waits to get
func RecieveFile(ctx context.Context, wg *sync.WaitGroup, rip net.IP){
	defer wg.Done()
	
	localLn := net.TCPAddr{
		IP: net.ParseIP("0.0.0.0"),
		Port: config.RecievePort,
	}

	ln, err := net.ListenTCP("tcp",&localLn)
	if err != nil {
		log.Fatalf("failed to create listener: %v", err)
	}else{
		log.Printf("waiting to recieve on port: %d", localLn.Port)
	}
	for { 
		conn, err := ln.Accept()
		if err != nil {
			// what are the possilbe errors?
			if errors.Is(err, net.ErrClosed){
				return
			}
			continue
		}
		log.Printf("Recieving on port: %d", localLn.Port)
		go downloadFile(conn)
	}
}

func downloadFile(conn net.Conn){
	log.Printf("Downloading file from %v....", conn.RemoteAddr().String())

	hdr, err := recieveHeader(conn)
	if err != nil{
		log.Printf("failed to recieve header: %v", err)
		return
	}
	
	log.Printf("head filename: %v", hdr.FileName)
	log.Printf("head checksum: %v", hdr.CheckSum)

	log.Printf("Continue Download? (yes/no): ")
	var resp string
	fmt.Scanln(&resp)
	if resp != "yes" {
		conn.Close()
		log.Print("Not downloading...")
		return
	}

	f, err := os.Create(hdr.FileName)
	if err != nil{
		log.Fatal("failed to create file: ", err)
	}
	defer f.Close()
	
	hash, bytesRecieved, err := shared.CopyAndHash(f,conn)
	if err != nil{
		log.Printf("error copying file: %v", err)
	}
	log.Printf("successfully copied %d bytes", bytesRecieved)
	log.Printf("hash of file is : %v", hash)

	if hash == hdr.CheckSum {
		log.Printf("hashes match")
	}

	log.Printf("finished downloading %v", hdr.FileName)
}


func recieveHeader(conn net.Conn) (models.Header, error) {

	var hdr models.Header
	var lenBuf [4]byte

	if _, err := io.ReadFull(conn, lenBuf[:]); err !=  nil {
		return models.Header{}, err
	}

	headerLen := binary.BigEndian.Uint32(lenBuf[:])
	log.Printf("header len: %v", headerLen)
	
	jsonData := make([]byte, headerLen)
	if _, err := io.ReadFull(conn, jsonData); err != nil {
		return models.Header{}, err
	}

	if err := json.Unmarshal(jsonData, &hdr); err != nil{
		return models.Header{}, err
	}

	return hdr, nil
}