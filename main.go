package main

import (
	"fmt"
	"os"

	"github.com/italobbarros/tcp-client-connect/internal/client"
	terminal "github.com/italobbarros/tcp-client-connect/internal/terminal"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ./tcpclient <address>")
		return
	}
	//signalCh := make(chan os.Signal, 1)
	endCh := make(chan struct{}, 1)
	//signal.Notify(signalCh, syscall.SIGTERM)

	myClient := client.NewClient(os.Args[1], endCh)
	gui := terminal.NewTerminal(myClient.ServerCommandCh, myClient.UserCommandCh, myClient.IsConnected)

	go gui.Create(endCh)
	go gui.ListenServerResponse(endCh)
	myClient.Start(gui.PrintStatus, gui.ClearInput)
}
