package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

var nodes map[string]net.Conn            //연결된 노드
var nodeMap map[string]map[string]string //통신으로 넘겨줄 데이터

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
	nodeMap = make(map[string]map[string]string)

	go handleServer()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("연결 승인 실패:", err)
			continue
		}
		//노드에게 nodes 전달
		nodeData, err := json.Marshal(nodeMap)
		if err != nil {
			fmt.Printf("nodeMap 직렬화 실패: %s\n", err)
			return
		}
		conn.Write(nodeData)

		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("노드 ListenPort 입력 오류:%s\n", err)
			return
		}
		listenPort := string(buffer[:bytesRead])
		dialPort := conn.RemoteAddr().String()

		nodes[listenPort] = conn
		nodeMap[dialPort] = map[string]string{
			"listenPort": listenPort,
			"dialPort":   dialPort,
		}
		fmt.Printf("새로운 노드 연결: %s\n", listenPort)

		go handleNode(conn)
	}
}

func handleServer() {
	for {
		reader := bufio.NewScanner(os.Stdin)
		fmt.Print("노드 포트 입력:")
		reader.Scan()
		port := reader.Text()
		port = strings.TrimSpace(port)

		_, err := strconv.Atoi(port)
		if err != nil {
			fmt.Printf("정수를 입력하시오.%s\n", err)
			continue
		}

		fmt.Print("메세지 입력:")
		reader.Scan()
		text := reader.Text()
		text = strings.TrimSpace(text)

		conn, ok := nodes["localhost:"+port]
		if ok {
			fmt.Printf("노드(%s)로 메세지 전달.", nodeMap[conn.RemoteAddr().String()]["listenPort"])
			conn.Write([]byte(text + "\n"))
		} else {
			fmt.Println("해당 포트의 노드는 존재하지 않습니다.")
		}
	}
}

func handleNode(conn net.Conn) {
	defer conn.Close()

	//노드 send에 의한 received
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			nodeDialAddress := conn.RemoteAddr().String()
			nodeAddress := nodeMap[nodeDialAddress]["listenPort"]
			fmt.Printf("노드 연결 끊김: %s\n", nodeAddress)
			delete(nodes, nodeDialAddress)
			break
		}

		fmt.Printf("Message Received(%s): %s", conn.RemoteAddr().String(), string(message))
	}
}
