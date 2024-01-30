package client

import (
	"fmt"
	"net"
	"time"
)

const (
	maxReconnectAttempts = 10
	reconnectInterval    = 2 * time.Second
)

// Client represents a TCP client
type Client struct {
	serverAddr        string
	ServerCommandCh   chan string
	UserCommandCh     chan string
	DoneCh            chan struct{}
	reconnectAttempts int
}

// NewClient initializes a new instance of the Client
func NewClient(serverAddr string) *Client {
	return &Client{
		serverAddr:        serverAddr,
		ServerCommandCh:   make(chan string),
		UserCommandCh:     make(chan string),
		DoneCh:            make(chan struct{}),
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
			case <-c.DoneCh:
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
		case <-c.DoneCh:
			return
		default:
			userInput := <-c.UserCommandCh
			conn.Write([]byte(userInput))
		}
	}
}

// Start initiates the client
func (c *Client) Start() {
	c.Connect()
}

func (c *Client) Stop() {
	close(c.DoneCh)
}
