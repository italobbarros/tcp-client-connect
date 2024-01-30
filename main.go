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
	serverAddr := os.Args[1]
	myClient := client.NewClient(serverAddr)
	gui := terminal.NewInterface(&myClient.ServerCommandCh, &myClient.UserCommandCh)
	go gui.Create()
	go gui.ListenServerResponse()
	myClient.Start()
}
