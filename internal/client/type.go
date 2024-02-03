package client

import (
	"net"
	"sync"
)

type Client struct {
	ClientId          int
	serverAddr        string
	InputData         chan DataType
	OutputData        chan DataType
	StatusCh          chan StatusMsg
	ReconnectCh       chan struct{}
	endCh             chan struct{}
	reconnectAttempts int
	conn              net.Conn
}

type ManagerClients struct {
	Map   map[int]*Client
	mutex sync.Mutex
}

type StatusMsg struct {
	Msg   string
	Color TextColors
}

type DataType struct {
	Data   []byte
	ConnId int
}

type TextColors int

const (
	TextRed TextColors = iota
	TextYellow
	TextGreen
	TextBlue
)
