package shared

import (
	"ftcli/config"
	"runtime"
)

// The following code was generate using ChatGPT
func PrintMemUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    config.Dlog.Printf("Alloc = %v MiB\n", m.Alloc/1024/1024)
    config.Dlog.Printf("TotalAlloc = %v MiB\n", m.TotalAlloc/1024/1024)
    config.Dlog.Printf("Sys = %v MiB\n", m.Sys/1024/1024)
    config.Dlog.Printf("NumGC = %v\n", m.NumGC)
}
