package main

import (
	"bufio"
	"bytes"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/47Cid/Castle/config"
	"github.com/47Cid/Castle/logger"
	"github.com/47Cid/Castle/message"
	"github.com/47Cid/Castle/pod"
)

func main() {
	// Load the configuration
	err := config.LoadConfig("./config/config.json")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Initialize the logger
	logger.InitProxyLogger()

	// Start listening for incoming connections
	listener, err := net.Listen("tcp", config.GetListenPort())
	if err != nil {
		logger.ProxyLog.Fatalf("Failed to start listener: %v", err)
	}

	logger.ProxyLog.Infof("Server listening on %s", listener.Addr().String())

	// WaitGroup for handling client connections
	var wg sync.WaitGroup

	// Handle incoming connections
	for {
		clientConn, err := listener.Accept()
		if err != nil {
			logger.ProxyLog.Errorf("Error accepting connection: %v", err)
			continue
		}

		wg.Add(1)
		go handleClient(&wg, clientConn)
	}

	// Wait for all client connections to be handled
	wg.Wait()
}

func handleClient(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()

	// Get the client IP
	clientIP := conn.RemoteAddr().(*net.TCPAddr).IP

	// Get the current timestamp
	timestamp := time.Now()

	// Read data from the client connection
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		logger.ProxyLog.Errorf("Error reading from client: %v", err)
		return
	}

	// Parse the HTTP request
	req, err := http.ReadRequest(bufio.NewReader(bytes.NewReader(buffer[:n])))
	if err != nil {
		logger.ProxyLog.Errorf("Error parsing HTTP request: %v", err)
		return
	}
	// Get the destination URL
	destinationURL := req.URL.String()

	msg := message.Message{
		Data:        buffer[:n],
		Destination: destinationURL,
		ClientIP:    clientIP,
		Timestamp:   timestamp,
	}

	// Call the verify function from waf_pod
	isValid := pod.VerifyMessage(msg)
	if !isValid {
		logger.ProxyLog.Errorf("Invalid request from client: %v", conn.RemoteAddr())
		return
	}

	// TODO: Handle the client request
}
