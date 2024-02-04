package tcp

import (
	"fmt"
	"time"
)

func NewManagerConnection(AddrList []string, endCh chan struct{}) *ManagerConnections {
	m := ManagerConnections{
		MapConnections:    make(map[int]*Connection),
		activeConnections: 0,
	}
	for i, addr := range AddrList {
		conn := NewConnection(i, addr, endCh, &m)
		m.AddConnections(i, conn)
		go conn.Start()
	}
	return &m
}

func (m *ManagerConnections) AddConnections(index int, client *Connection) {
	m.mutexConnections.Lock()
	m.MapConnections[index] = client
	m.mutexConnections.Unlock()
}

func (m *ManagerConnections) GetNumberConnections() int {
	m.mutexConnections.Lock()
	number := len(m.MapConnections)
	m.mutexConnections.Unlock()
	return number
}

func (m *ManagerConnections) SendDataToConnections(data DataType) []int {
	m.mutexConnections.Lock()
	var clientIds []int
	for _, client := range m.MapConnections {
		data.ConnId = client.Id
		clientIds = append(clientIds, client.Id)
		if !client.IsConnected() {
			continue
		}
		client.OutputData <- data
	}
	m.mutexConnections.Unlock()
	return clientIds
}

func (m *ManagerConnections) ReceiveDataToConnections(received chan DataType, StatusCh chan StatusMsg, StatusInfoCh chan StatusMsg) {
	m.mutexConnections.Lock()
	m.StatusInfoCh = StatusInfoCh
	for _, client := range m.MapConnections {
		client.StatusCh = StatusCh
		go func(c *Connection) {
			for {
				select {
				case <-c.endCh:
					return
				case v := <-c.InputData:
					received <- v
				}
			}
		}(client)
	}
	m.mutexConnections.Unlock()
}

func (m *ManagerConnections) AddActiveConnections() {
	m.mutexConnections.Lock()
	m.activeConnections++
	m.mutexConnections.Unlock()
	go m.PrintStatusInfo()
}
func (m *ManagerConnections) RemoveActiveConnections() {
	m.mutexConnections.Lock()
	m.activeConnections--
	m.mutexConnections.Unlock()
	go m.PrintStatusInfo()
}

func (m *ManagerConnections) Start(endCh chan struct{}) {
	select {
	case <-endCh:
		return
	}
}

func (m *ManagerConnections) PrintStatusInfo() {
	time.Sleep(time.Millisecond * 1)
	total := m.GetNumberConnections()
	actives := m.activeConnections
	deactives := total - actives
	color := TextBlue

	if actives == 0 {
		color = TextRed
	}
	if actives == total {
		color = TextGreen
	}
	m.StatusInfoCh <- StatusMsg{
		Msg:   fmt.Sprintf("T:%d | C:%d | D:%d", total, actives, deactives),
		Color: color,
	}
}
