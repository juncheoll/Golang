package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
)

var nodes map[string]net.Conn //연결된 노드
var nodeMap map[string]string

func main() {
	serverAddress := "localhost:8080"
	runHeadServer(serverAddress)
}

func runHeadServer(serverAddress string) {
	ln, err := net.Listen("tcp", serverAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	fmt.Println("총괄 서버 시작 - ", serverAddress)
	nodes = make(map[string]net.Conn)
	nodeMap = make(map[string]string)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("연결 승인 실패:", err)
			continue
		}
		nodes[conn.RemoteAddr().String()] = conn
		nodeMap[fmt.Sprintf("node%d", len(nodes))] = conn.RemoteAddr().String()
		fmt.Printf("새로운 노드 연결: %s\n", nodeMap[fmt.Sprintf("node%d", len(nodes))])

		go handleNode(conn)
	}
}

func handleNode(conn net.Conn) {
	defer conn.Close()

	//노드에게 nodes 전달
	nodeData, err := json.Marshal(nodeMap)
	if err != nil {
		fmt.Printf("nodeMap 직렬화 실패: %s\n", err)
		return
	}
	conn.Write(nodeData)

	//노드 send에 의한 received
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("노드 연결 끊김: %s\n", conn.RemoteAddr().String())
			delete(nodes, conn.RemoteAddr().String())
			break
		}

		fmt.Printf("Message Received(%s): %s", conn.RemoteAddr().String(), string(message))
	}
}
