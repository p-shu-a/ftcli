package main

import (
	"context"
	"flag"
	"ftcli/internal/receive"
	"ftcli/internal/send"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

func main(){
	// ideally, you can just do something like "--sender" or "--recieve" without the '-role' tag
	role := flag.String("role", "" , "Designate yourself as sender (send) or reciever (recv)")
	rip := flag.String("rip", "127.0.0.1", "ip address of the peer")
	file := flag.String("file", "", "file to transfer")
	//pass := flag.String("pass", "", "password used to encrypt file")
	flag.Parse()

	// do validation of flags here
	// Validate password. user must enter password.
	// accomodate no password, infrom user of implications
	// if *pass == "" {
	// 	log.Fatal("must specify a password for transfer")
	// }


	// Validate provided role as one of two valid roles
	var userRole string

	switch *role {
	case "send":
		userRole = "send"
	case "recv":
		userRole = "recv"
	default:
		log.Fatal("unrecognized role. use send or recv")
	}


	// Validate provided IP as legit IP...
	peerAddr := net.ParseIP(*rip)
	if peerAddr == nil{
		log.Fatal("Must enter valid IP")
	}

	var wg sync.WaitGroup

	// Validate provided file. 
	// if userRole is not send, then silently ignore
	if userRole == "send" {
		// go can't expand ~, so you manully have to replace the ~ with the userdir
		home, _ := os.UserHomeDir()
		path := strings.Replace(*file, "~", home, 1)

		file, err := os.Open(path)
		if err != nil{
			log.Fatalf("fail failure: %v", err)
		}
		wg.Add(1)
		go send.SendFile(context.TODO(), &wg, file, peerAddr)
	}else{
		wg.Add(1)
		go receive.ReceiveFile(context.TODO(), &wg)
	}

	wg.Wait()
	log.Print("operations complete. good-bye 👋")
}


