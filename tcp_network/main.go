package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var clients []net.Conn

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("please provide the argument")
		return
	}

	port := args[1]

	mainAddress := "127.0.0.1:8080"
	nodeAddress := "127.0.0.1:" + port

	go startServer(nodeAddress)
	go startClient(mainAddress)
}

func startClient(address string) {
	connClient, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println("error:", err)
	}

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Text to send: (EXIT : exit)")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "exit" {
			connClient.Close()
			fmt.Println("종료")
			return
		}

		fmt.Fprint(connClient, text+"\n")
	}
}

func startServer(serverAddress string) {
	fmt.Println("Starting...")

	ln, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Printf("New client connected: %s\n", conn.RemoteAddr().String())
		clients = append(clients, conn)

		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", conn.RemoteAddr().String())

			for idx := 0; idx < len(clients); idx++ {
				if clients[idx] == conn {
					clients = append(clients[:idx], clients[idx+1:]...)
					break
				}
			}

			break
		}

		fmt.Print("Message Received:" + string(message))

		conn.Write([]byte(message + "\n"))
	}
}
