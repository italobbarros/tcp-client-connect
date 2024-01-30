package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		// Aguarde um sinal para encerrar a aplicação
		<-signalCh

		// Feche os canais e encerre a aplicação
		close(myClient.ServerCommandCh)
		close(myClient.UserCommandCh)
		myClient.Stop()

		os.Exit(0)
	}()
	go gui.Create(myClient.DoneCh)
	go gui.ListenServerResponse(myClient.DoneCh)
	myClient.Start()
}
