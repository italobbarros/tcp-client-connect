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
	// Verifica se a ajuda foi solicitada
	endCh := make(chan struct{}, 1)
	managerClients := client.ManagerClients{
		Map: make(map[int]*client.Client),
	}

	for i, addr := range config.Addr {
		myClient := client.NewClient(i, addr, endCh)
		managerClients.AddClients(i, myClient)
	}

	gui := terminal.NewTerminal(&managerClients)
	managerClients.ReceiveDataToClients(gui.Input, gui.StatusCh)
	go gui.Create(endCh)
	go gui.ListenServerResponse(endCh)
	managerClients.Start(endCh)
}
