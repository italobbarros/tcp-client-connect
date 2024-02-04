package client

func (m *ManagerConnections) AddConnections(index int, client *Connection) {
	m.mutex.Lock()
	m.Map[index] = client
	m.mutex.Unlock()
}

func (m *ManagerConnections) GetNumberConnections() int {
	m.mutex.Lock()
	number := len(m.Map)
	m.mutex.Unlock()
	return number
}

func (m *ManagerConnections) SendDataToConnections(data DataType) []int {
	m.mutex.Lock()
	var clientIds []int
	for _, client := range m.Map {
		data.ConnId = client.Id
		clientIds = append(clientIds, client.Id)
		client.OutputData <- data
	}
	m.mutex.Unlock()
	return clientIds
}

func (m *ManagerConnections) ReceiveDataToConnections(received chan DataType, StatusCh chan StatusMsg) {
	m.mutex.Lock()
	for _, client := range m.Map {
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
	m.mutex.Unlock()
}

func (m *ManagerConnections) Start(endCh chan struct{}) {
	m.mutex.Lock()
	for _, client := range m.Map {
		go client.Start()
	}
	m.mutex.Unlock()
	select {
	case <-endCh:
		return
	}
}
