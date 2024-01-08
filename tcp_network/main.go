package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("please provide the argument")
		return
	}

	connType := args[1]

	address := "127.0.0.1:8081"

	if connType == "server" {
		startServer()
	} else if connType == "client" {
		startClient(address)
	} else {
		log.Println("please provide the argument (server or client)")
	}

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

func startServer() {
	fmt.Println("Starting...")

	ln, err := net.Listen("tcp", ":8081")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		connServer, err := ln.Accept()
		if err != nil {
			log.Println("Failed to accept connection:", err)
			continue
		}
		fmt.Printf("New client connected: %s\n", connServer.RemoteAddr().String())

		go handleClient(connServer)
	}
}

func handleClient(connServer net.Conn) {
	defer connServer.Close()

	for {
		message, err := bufio.NewReader(connServer).ReadString('\n')
		if err != nil {
			fmt.Printf("Client disconnected: %s\n", connServer.RemoteAddr().String())
			break
		}

		fmt.Print("Message Received:" + string(message))

		connServer.Write([]byte(message + "\n"))
	}
}
