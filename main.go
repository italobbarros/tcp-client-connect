package main

import (
	"fmt"

	"github.com/italobbarros/tcp-client-connect/internal/arg"
	"github.com/italobbarros/tcp-client-connect/internal/tcp"
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
	managerConnections := tcp.NewManagerConnection(config.Addr, endCh)

	gui := terminal.NewTerminal(managerConnections)
	managerConnections.ReceiveDataToConnections(gui.Input, gui.StatusCh, gui.StatusInfoCh)
	go gui.Create(endCh)
	go gui.ListenServerResponse(endCh)
	managerConnections.Start(endCh)

}
