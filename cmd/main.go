package main

import (
	"context"
	"flag"
	"ftcli/internal/receive"
	"ftcli/internal/send"
	"ftcli/internal/shared"
	"log"
	"net"
	_ "net/http/pprof"
	"os"
	"strings"
	"sync"
)

func main() {

	// run pprof
	go shared.RunPProf()

	// ideally, you can just do something like "--sender" or "--recieve" without the '-role' tag
	role := flag.String("role", "", "Designate yourself as sender (send) or receiver (recv)")
	rip := flag.String("rip", "127.0.0.1", "ip address of peer")
	file := flag.String("file", "", "file to transfer")
	pass := flag.String("pass", "", "password used to encrypt file")
	flag.Parse()

	// Validate password. user must enter password
	if *pass == "" {
		panic("must specify a password for transfer")
	}

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
	if peerAddr == nil {
		log.Fatal("must enter valid IP")
	}

	var wg sync.WaitGroup

	// Validate provided file.
	// if userRole is not send, then silently ignore
	if userRole == "send" {
		// go can't expand ~, so you manully have to replace the ~ with the userdir
		home, _ := os.UserHomeDir()
		path := strings.Replace(*file, "~", home, 1)

		// open the file you're going to send
		file, err := os.Open(path)
		if err != nil {
			log.Fatalf("failed to open file:  %v", err)
		}

		wg.Add(1)
		if err := send.SendFile(context.TODO(), &wg, file, peerAddr, *pass); err != nil {
			log.Fatalf("failed to send file: %v", err)
		}
	} else {

		wg.Add(1)
		if err := receive.ReceiveFile(context.TODO(), &wg, *pass); err != nil{
			log.Fatal(err)
		}
	}

	wg.Wait()
	log.Print("operations complete. good-bye 👋")
}
