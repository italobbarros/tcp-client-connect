package main

import (
	"fmt"

	"github.com/italobbarros/tcp-client-connect/internal/arg"
	"github.com/italobbarros/tcp-client-connect/internal/client"
	terminal "github.com/italobbarros/tcp-client-connect/internal/terminal"
)

func main() {
	config, err := arg.ParseFlags()
	if err != nil {
		errStr := err.Error()
		if errStr != "" {
			fmt.Println("Error: " + errStr)
		}
		return
	}
	//fmt.Println(config)
	// Verifica se a ajuda foi solicitada
	endCh := make(chan struct{}, 1)
	managerClients := client.ManagerConnections{
		Map: make(map[int]*client.Connection),
	}

	for i, addr := range config.Addr {
		myClient := client.NewConnection(i, addr, endCh)
		managerClients.AddConnections(i, myClient)
	}

	gui := terminal.NewTerminal(&managerClients)
	managerClients.ReceiveDataToConnections(gui.Input, gui.StatusCh)
	go gui.Create(endCh)
	go gui.ListenServerResponse(endCh)
	managerClients.Start(endCh)

}
