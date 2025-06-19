package config

import (
	"log"
	"os"
)

// Shared configs live here

// Port number on which the receving peer listens
const ReceivePort = 7891

// File chunk size, in bytes, for sending encrypted files in chunks
const FileChunkSize uint = 4096

// Cost settings for argon2. These are standard values.
var (
	Time      uint32 = 3
	Memory    uint32 = 32 * 1024
	Threads   uint8  = 4
	KeyLength uint32 = 32
)

// Debug logger. set out to io.discard to stop logging
var Dlog = log.New(os.Stdout, "DEBUG: ", log.Lmsgprefix)
var Slog = log.New(os.Stdout, "STATs: ", log.Lmsgprefix)
