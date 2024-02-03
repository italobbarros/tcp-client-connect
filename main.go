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
	myClient := client.NewClient(config.Addr, endCh)
	gui := terminal.NewTerminal(myClient.ServerCommandCh, myClient.UserCommandCh, myClient.IsConnected)

	go gui.Create(endCh)
	go gui.ListenServerResponse(endCh)
	myClient.Start(gui.PrintStatus, gui.ClearAll)
}
