package shared

import (
	"ftcli/config"
	"log"
	"net/http"
	"runtime"
)

// The following function was generated using ChatGPT
// Mem analysis should be improved to record memstats, not print out to stdout.
// shoudl be able to query when needed to see usage at some point
// hard to get "peak" mem usages otherwise
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	config.Slog.Printf("Alloc (peak) = %v MiB\n", m.Alloc/1024/1024)
	config.Slog.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
	config.Slog.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
	config.Slog.Printf("NumGC = %v\n", m.NumGC)
}

func RunPProf() {
	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		log.Printf("failed to run pprof: %v", err)
	}
	log.Println("pprof running on port 6060")
}
