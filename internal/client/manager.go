package client

func (m *ManagerClients) AddClients(index int, client *Client) {
	m.mutex.Lock()
	m.Map[index] = client
	m.mutex.Unlock()
}

func (m *ManagerClients) GetNumberClients() int {
	m.mutex.Lock()
	number := len(m.Map)
	m.mutex.Unlock()
	return number
}

func (m *ManagerClients) SendDataToClients(data DataType) []int {
	m.mutex.Lock()
	var clientIds []int
	for _, client := range m.Map {
		data.ConnId = client.ClientId
		clientIds = append(clientIds, client.ClientId)
		client.OutputData <- data
	}
	m.mutex.Unlock()
	return clientIds
}

func (m *ManagerClients) ReceiveDataToClients(received chan DataType, StatusCh chan StatusMsg) {
	m.mutex.Lock()
	for _, client := range m.Map {
		client.StatusCh = StatusCh
		go func(c *Client) {
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

func (m *ManagerClients) Start(endCh chan struct{}) {
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
