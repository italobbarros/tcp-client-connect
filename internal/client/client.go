package client

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxReconnectAttempts = 10
	reconnectInterval    = 2 * time.Second
)

// Client represents a TCP client
type Client struct {
	serverAddr        string
	stopCh            chan os.Signal
	ServerCommandCh   chan string
	UserCommandCh     chan string
	doneCh            chan struct{}
	reconnectAttempts int
}

// NewClient initializes a new instance of the Client
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr:        serverAddr,
		stopCh:            make(chan os.Signal, 1),
		ServerCommandCh:   make(chan string),
		UserCommandCh:     make(chan string),
		doneCh:            make(chan struct{}),
		reconnectAttempts: 0,
	}
}

// Connect establishes a connection to the server
func (c *Client) Connect() {
	var conn net.Conn
	var err error

	// Attempt to reconnect until reaching the maximum number of attempts
	for c.reconnectAttempts < maxReconnectAttempts {
		conn, err = net.Dial("tcp", c.serverAddr)
		if err == nil {
			break
		}

		fmt.Printf("Error connecting (attempt %d/%d): %v. Retrying in %v...\n",
			c.reconnectAttempts+1, maxReconnectAttempts, err, reconnectInterval)

		c.reconnectAttempts++
		time.Sleep(reconnectInterval)
	}

	// If unable to connect after the maximum number of attempts, terminate the application
	if err != nil {
		fmt.Printf("Unable to reconnect after %d attempts. Exiting the application.\n", maxReconnectAttempts)
	}

	// Reset the retry counter after a successful connection
	c.reconnectAttempts = 0

	defer conn.Close()

	// Receive and print server commands
	go func() {
		for {
			select {
			case <-c.doneCh:
				return
			default:
				buffer := make([]byte, 4096)
				_, err := conn.Read(buffer)
				if err != nil {
					fmt.Println("Error receiving response:", err)
					return
				}
				// Send server command to the channel
				c.ServerCommandCh <- string(buffer)
			}
		}
	}()

	// Receive user input and send commands to the server
	for {
		select {
		case <-c.doneCh:
			return
		default:
			// Receive user input
			//scanner := bufio.NewScanner(os.Stdin)
			//fmt.Print("Enter a command: ")
			//if !scanner.Scan() {
			//	return
			//}
			//userInput := scanner.Text()
			// Send user command to the server
			userInput := <-c.UserCommandCh
			conn.Write([]byte(userInput))
		}
	}
}

// Start initiates the client
func (c *Client) Start() {
	signal.Notify(c.stopCh, syscall.SIGINT, syscall.SIGTERM)
	go c.Connect()

	// Wait for signals to stop the program
	<-c.stopCh
	close(c.ServerCommandCh)
	close(c.UserCommandCh)
}

func (c *Client) Stop() {
	close(c.doneCh)
}
