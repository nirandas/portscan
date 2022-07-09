package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"
)

var host = flag.String("h", "", "Host domain or ip to scan")
var minPort = flag.Int("min", 1, "The starting port number for scan")
var maxPort = flag.Int("max", 1024, "Max port number to stop scan")
var concurrency = flag.Int("c", 10, "The concurrent connections to open")

func main() {
	parseFlags()
	fmt.Printf("Scanning %s port %d to %d with %d concurrent connections\r\n", *host, *minPort, *maxPort, *concurrency)

	portChan := make(chan int, 1)
	var ports []int
	var lock sync.Mutex
	var wg sync.WaitGroup
	start := time.Now()

	for i := 0; i < *concurrency; i++ {
		go func() {
			for {
				port := <-portChan
				if attempt(port) {
					lock.Lock()
					ports = append(ports, port)
					lock.Unlock()
				}
				wg.Done()
			}
		}()
	}

	for i := *minPort; i <= *maxPort; i++ {
		wg.Add(1)
		portChan <- i
	}

	wg.Wait()
	fmt.Println("Done in", time.Now().Sub(start))
	fmt.Println("Open ports", ports)
}

func attempt(port int) bool {
	defer func() {
		recover()
	}()

	address := fmt.Sprintf("%s:%d", *host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}

func parseFlags() {
	flag.Parse()
	if *host == "" {
		fmt.Println("Host is required")
		os.Exit(0)
	}
	if *minPort < 1 || *maxPort > 65535 || *minPort > *maxPort {
		fmt.Println("Min port should be > 0, max port should be < 65536 and min port should not be > than the max port")
		os.Exit(0)
	}
	if *concurrency < 1 || *concurrency > 100 {
		fmt.Println("Concurrent connections should be between 1 and 100 inclusive")
		os.Exit(0)
	}
}
