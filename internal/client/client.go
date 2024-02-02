package client

import (
	"fmt"
	"net"
	"time"
)

const (
	maxReconnectAttempts = 100
	reconnectInterval    = 2 * time.Second
)

// Client represents a TCP client
type Client struct {
	serverAddr        string
	ServerCommandCh   chan string
	UserCommandCh     chan string
	DoneCh            chan struct{}
	reconnectAttempts int
	conn              net.Conn
	Print             func(value string)
	Clear             func()
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

// Start initiates the client
func (c *Client) Connect() {
	var err error
	for {
		c.conn, err = net.Dial("tcp", c.serverAddr)
		if err == nil {
			break
		}

		msg := fmt.Sprintf("Error connecting (attempt %d): %v. Retrying in %v...\n",
			c.reconnectAttempts+1, err, reconnectInterval)
		c.Print(msg)
		c.reconnectAttempts++
		time.Sleep(reconnectInterval)
	}

	// If unable to connect after the maximum number of attempts, terminate the application
	if err != nil {
		print(fmt.Sprintf("Unable to reconnect after %d attempts. Exiting the application.\n", maxReconnectAttempts))
	}

	// Reset the retry counter after a successful connection
	c.reconnectAttempts = 0
	c.Clear()
	c.Print("Conectou!")
}

func (c *Client) Start(print func(value string), clear func()) {
	time.Sleep(100 * time.Millisecond)
	c.Print = print
	c.Clear = clear
	c.Connect()
	defer c.conn.Close()
	// Receive and print server commands
	go c.startRead()
	go c.startWrite()
	select {}
}

func (c *Client) Stop() {
	close(c.DoneCh)
}

func (c *Client) startRead() {
	for {
		select {
		case <-c.DoneCh:
			c.DoneCh = make(chan struct{})
			c.Connect()
			continue
		default:
			if c.reconnectAttempts != 0 {
				continue
			}
			buffer := make([]byte, 4096)
			_, err := c.conn.Read(buffer)
			if err != nil {
				c.Print(fmt.Sprintf("Error receiving response:%s", err.Error()))
				continue
			}
			// Send server command to the channel
			c.ServerCommandCh <- string(buffer)
		}
	}
}

func (c *Client) startWrite() {
	for {
		select {
		case <-c.DoneCh:
			continue
		default:
			if c.reconnectAttempts != 0 {
				continue
			}
			c.conn.Write([]byte(<-c.UserCommandCh))
		}
	}
}
