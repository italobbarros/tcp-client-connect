package tcp

import (
	"net"
	"sync"
)

type Connection struct {
	Id                int
	serverAddr        string
	InputData         chan DataType
	OutputData        chan DataType
	StatusCh          chan StatusMsg
	ReconnectCh       chan struct{}
	endCh             chan struct{}
	reconnectAttempts int
	conn              net.Conn
	manager           *ManagerConnections
}

type ManagerConnections struct {
	StatusInfoCh      chan StatusMsg
	MapConnections    map[int]*Connection
	mutexConnections  sync.Mutex
	activeConnections int
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
