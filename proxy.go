package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	// Create a channel to receive signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// WaitGroup for handling client connections
	var wg sync.WaitGroup

	go func() {
		<-signalChan
		logger.ProxyLog.Info("Received termination signal. Initiating graceful exit.")
		os.Exit(0)
	}()

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
}

func forwardMessage(conn net.Conn, msg message.Message) {
	logger.ProxyLog.Infof("Valid request from client")
	targetConn, err := net.Dial("tcp", config.GetRemoteServerAddr())
	if err != nil {
		log.Printf("Error connecting to target server: %v", err)
		return
	}
	defer targetConn.Close()
	go func() {
		_, err = io.Copy(conn, targetConn)
		if err != nil {
			log.Printf("Error copying data from server to client: %v", err)
		}
	}()
	_, err = io.Copy(targetConn, conn)
	if err != nil {
		log.Printf("Error copying data from client to server: %v", err)
	}
}

func handleClient(wg *sync.WaitGroup, conn net.Conn) {
	defer wg.Done()
	defer conn.Close()

	// Get the client IP
	clientIP := conn.RemoteAddr().(*net.TCPAddr).IP

	// Get the current timestamp
	timestamp := time.Now()

	msg := message.Message{
		ClientIP:  clientIP,
		Timestamp: timestamp,
	}

	// Log the message
	logger.ProxyLog.Infof("Received message: %+v", msg)

	// Call the verify function from waf_pod
	isValid := pod.VerifyMessage(msg)
	if !isValid {
		logger.ProxyLog.Errorf("Invalid request from client: %v", conn.RemoteAddr())

		// Define the path to the HTML file
		htmlFilePath := "assets/error.html"

		// Read the HTML file
		htmlData, err := os.ReadFile(htmlFilePath)
		if err != nil {
			logger.ProxyLog.Errorf("Error reading HTML file: %v", err)
			return
		}

		// Write the HTTP response status line
		_, err = conn.Write([]byte("HTTP/1.1 403 Forbidden\r\n"))
		if err != nil {
			logger.ProxyLog.Errorf("Error writing response: %v", err)
			return
		}

		// Write the Content-Type header
		_, err = conn.Write([]byte("Content-Type: text/html\r\n\r\n"))
		if err != nil {
			logger.ProxyLog.Errorf("Error writing response: %v", err)
			return
		}

		// Write the HTML data
		_, err = conn.Write(htmlData)
		if err != nil {
			logger.ProxyLog.Errorf("Error writing response: %v", err)
		}
		return
	}
	//Forward the message to the server
	forwardMessage(conn, msg)
}
