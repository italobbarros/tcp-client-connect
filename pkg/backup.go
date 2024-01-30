package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	maxReconnectAttempts = 10
	reconnectInterval    = 2 * time.Second
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Uso: ./tcpclient <endereço>")
		return
	}

	serverAddr := os.Args[1]

	// Capturando sinais para interromper o programa
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM)

	// Canal para receber comandos do servidor
	serverCommandCh := make(chan string)

	// Canal para enviar comandos do usuário ao servidor
	userCommandCh := make(chan string)

	// Canal para sinalizar a goroutine principal para encerrar
	doneCh := make(chan struct{})

	// Contador de tentativas de reconexão
	reconnectAttempts := 0

	// Iniciando as requisições concorrentes
	go func() {
		defer close(doneCh)

		var conn net.Conn
		var err error

		// Tentar reconectar até atingir o número máximo de tentativas
		for reconnectAttempts < maxReconnectAttempts {
			conn, err = net.Dial("tcp", serverAddr)
			if err == nil {
				break
			}

			fmt.Printf("Erro ao conectar (tentativa %d/%d): %v. Tentando novamente em %v...\n",
				reconnectAttempts+1, maxReconnectAttempts, err, reconnectInterval)

			reconnectAttempts++
			time.Sleep(reconnectInterval)
		}

		// Se não foi possível conectar após o número máximo de tentativas, encerrar a aplicação
		if err != nil {
			fmt.Printf("Não foi possível reconectar após %d tentativas. Encerrando a aplicação.\n", maxReconnectAttempts)
		}

		// Resetar o contador de tentativas após uma conexão bem-sucedida
		reconnectAttempts = 0

		defer conn.Close()

		// Receber e imprimir comandos do servidor
		go func() {
			for {
				select {
				case <-doneCh:
					return
				default:
					buffer := make([]byte, 4096)
					_, err := conn.Read(buffer)
					if err != nil {
						fmt.Println("Erro ao receber resposta:", err)
						return
					}
					// Enviar comando do servidor para o canal
					//serverCommandCh <- string(buffer)
					fmt.Printf("Recebeu:%s\n", string(buffer))

				}
			}
		}()

		// Receber input do usuário e enviar comandos ao servidor
		for {
			select {
			case <-doneCh:
				return
			default:
				// Receber input do usuário
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Print("Digite um comando: ")
				if !scanner.Scan() {
					return
				}
				userInput := scanner.Text()
				// Enviar comando do usuário ao servidor
				conn.Write([]byte(userInput))
			}
		}

	}()

	// Aguardando sinais para interromper o programa
	<-stopCh
	close(serverCommandCh)
	close(userCommandCh)

	// Esperar pela conclusão da goroutine principal
	<-doneCh
}
