package config

import (
	"log"
	"os"
)

// Shared configs live here

var ReceivePort int = 7891

// Cost settings for argon2
var (
	Time      uint32 = 3
	Memory    uint32 = 32 * 1024
	Threads   uint8  = 4
	KeyLength uint32 = 32
)

var Dlog = log.New(os.Stdout, "DEBUG: ", log.Lmsgprefix)
//var debug = true

func EnableDebugLog(debug bool){
	if debug {
		Dlog.SetOutput(os.Stdout)
	}
}