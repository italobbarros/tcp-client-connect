package tcp

import (
	"fmt"
	"net"
	"time"
)

const (
	maxReconnectAttempts = 100
)

var reconnectInterval = 1 * time.Second

// Connection represents a TCP client

// NewConnection initializes a new instance of the Connection
func NewConnection(i int, serverAddr string, endCh chan struct{}, manager *ManagerConnections) *Connection {
	return &Connection{
		Id:                i,
		serverAddr:        serverAddr,
		InputData:         make(chan DataType),
		OutputData:        make(chan DataType),
		ReconnectCh:       make(chan struct{}),
		endCh:             endCh,
		reconnectAttempts: 0,
		manager:           manager,
	}
}

// Start initiates the client
func (c *Connection) Connect() {
	var err error
	for {
		select {
		case <-c.endCh:
			return
		case <-time.After(reconnectInterval):
			//fmt.Println(c.serverAddr)
			c.PrintStatus("Conectando...", TextBlue)
			c.conn, err = net.Dial("tcp", c.serverAddr)

			msg := fmt.Sprintf("Error connecting (attempt %d): Retrying in %v...",
				c.reconnectAttempts, reconnectInterval)
			c.PrintStatus(msg, TextRed)
			c.reconnectAttempts++
			reconnectInterval = 5 * time.Second
			if err != nil {
				continue
			}

			// Reset the retry counter after a successful connection
			c.reconnectAttempts = 0
			c.PrintStatus("Conectado", TextGreen)
			c.manager.AddActiveConnections()
			go c.startRead()
			go c.startWrite()
			return
		}
	}
}

func (c *Connection) Start() {
	//fmt.Println("Starting clientid ", c.Id)
	go c.Connect()
	go func() {
		<-c.endCh
		defer c.Stop()
	}()
}

func (c *Connection) Stop() {
	close(c.InputData)
	close(c.OutputData)
	close(c.ReconnectCh)

}

func (c *Connection) startRead() {
	for {
		select {
		case <-c.endCh:
			return
		default:
			if !c.IsConnected() {
				continue
			}
			if c.conn == nil {
				continue
			}
			buffer := make([]byte, 4096)
			_, err := c.conn.Read(buffer)
			if err != nil {
				c.PrintStatus(err.Error(), TextRed)
				c.ReconnectCh = make(chan struct{})
				c.conn.Close()
				c.PrintStatus("Desconectado!", TextRed)
				c.manager.RemoveActiveConnections()
				go c.Connect()
				return
			}
			// Send server command to the channel
			if c.InputData != nil {
				c.InputData <- DataType{
					Data:   buffer,
					ConnId: c.Id,
				}
			}
		}
	}
}

func (c *Connection) IsConnected() bool {
	return c.reconnectAttempts == 0
}

func (c *Connection) startWrite() {
	for {
		select {
		case <-c.endCh:
			return
		default:
			if c.reconnectAttempts != 0 {
				continue
			}
			if c.conn == nil {
				continue
			}
			data := <-c.OutputData
			if !c.IsConnected() {
				continue
			}
			_, err := c.conn.Write(data.Data)
			if err != nil {
				return
			}
		}
	}
}

func (c *Connection) PrintStatus(value string, color TextColors) {
	c.StatusCh <- StatusMsg{
		Msg:   fmt.Sprintf("\n[%s](%s - Conn%d): %s", time.Now().Format("2006-01-02 15:04:05"), c.serverAddr, c.Id, value),
		Color: color,
	}
}
