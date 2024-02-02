package client

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/italobbarros/tcp-client-connect/internal/terminal"
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
	ReconnectCh       chan struct{}
	endCh             chan struct{}
	reconnectAttempts int
	conn              net.Conn
	PrintStatus       func(value string, color terminal.TeminalColors)
	Clear             func()
}

// NewClient initializes a new instance of the Client
func NewClient(serverAddr string, endCh chan struct{}) *Client {
	return &Client{
		serverAddr:        serverAddr,
		ServerCommandCh:   make(chan string),
		UserCommandCh:     make(chan string),
		ReconnectCh:       make(chan struct{}),
		endCh:             endCh,
		reconnectAttempts: 0,
	}
}

// Start initiates the client
func (c *Client) Connect() {
	var err error
	for {
		select {
		case <-c.endCh:
			return
		case <-time.After(reconnectInterval):
			c.PrintStatus(c.serverAddr+" -> Conectando...", terminal.Blue)
			c.conn, err = net.Dial("tcp", c.serverAddr)

			msg := fmt.Sprintf(c.serverAddr+" -> Error connecting (attempt %d): Retrying in %v...",
				c.reconnectAttempts, reconnectInterval)
			c.PrintStatus(msg, terminal.Red)
			c.reconnectAttempts++
			if err != nil {
				continue
			}

			// Reset the retry counter after a successful connection
			c.reconnectAttempts = 0
			c.Clear()
			c.PrintStatus(c.serverAddr+" -> Conectado", terminal.Green)
			return
		}
	}
}

func (c *Client) Start(printStatus func(value string, color terminal.TeminalColors), clear func()) {
	c.PrintStatus = printStatus
	c.Clear = clear
	go c.Connect()
	// Receive and print server commands
	go c.startRead()
	go c.startWrite()
	select {
	case <-c.endCh:
		close(c.ServerCommandCh)
		close(c.UserCommandCh)
		c.Stop()
		os.Exit(0)
	}
}

func (c *Client) Stop() {
	close(c.ReconnectCh)
}

func (c *Client) startRead() {
	for {
		select {
		case <-c.endCh:
			return
		case <-c.ReconnectCh:
			c.ReconnectCh = make(chan struct{})
			c.conn.Close()
			c.PrintStatus(c.serverAddr+" -> Desconectado!", terminal.Red)
			c.Connect()
			continue
		default:
			if c.reconnectAttempts != 0 {
				continue
			}
			if c.conn == nil {
				continue
			}
			buffer := make([]byte, 4096)
			_, err := c.conn.Read(buffer)
			if err != nil {
				c.PrintStatus(c.serverAddr+" -> "+err.Error(), terminal.Red)
				c.Stop()
				continue
			}
			// Send server command to the channel
			if c.ServerCommandCh != nil {
				c.ServerCommandCh <- string(buffer)
			}
		}
	}
}

func (c *Client) IsConnected() bool {
	return c.reconnectAttempts == 0
}

func (c *Client) startWrite() {
	for {
		select {
		case <-c.endCh:
			return
		case <-c.ReconnectCh:
			continue
		default:
			if c.reconnectAttempts != 0 {
				continue
			}
			if c.conn == nil {
				continue
			}
			v := <-c.UserCommandCh
			_, err := c.conn.Write([]byte(v))
			if err != nil {
				c.PrintStatus(err.Error(), terminal.Red)
				c.Stop()
			}
		}
	}
}
