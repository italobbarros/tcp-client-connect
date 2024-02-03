package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal("Erro ao iniciar o servidor:", err)
	}

	defer listener.Close()

	fmt.Println("Servidor ouvindo em :8081")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Erro ao aceitar a conexão:", err)
			continue
		}

		fmt.Printf("Conexão recebida de %s\n", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Erro ao ler dados da conexão %s: %v\n", conn.RemoteAddr(), err)
			return
		}

		data := buffer[:n]
		fmt.Printf("Dados recebidos de %s: %s\n", conn.RemoteAddr(), data)
		n, err = conn.Write(data)
		if err != nil {
			log.Printf("Erro ao ler dados da conexão %s: %v\n", conn.RemoteAddr(), err)
			return
		}
	}
}
