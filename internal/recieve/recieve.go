package recieve

import (
	"context"
	"errors"
	"ftcli/config"
	"ftcli/internal/shared"
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
	f, err := os.Create("tmp.file")
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

	log.Printf("finish download form %v....", conn.RemoteAddr().String())
}