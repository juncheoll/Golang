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
	"sync"
)

var nodes map[string]net.Conn

func main() {
	nodes = make(map[string]net.Conn)

	args := os.Args

	if len(args) < 2 {
		fmt.Println("node 포트를 주세요")
		return
	}
	port := args[1]
	_, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println("node 포트를 정수로 입력해주세요")
	}

	nodeAddress := "localhost:" + port
	headAddress := "localhost:8080"

	ln, err := net.Listen("tcp", nodeAddress)
	fmt.Println(nodeAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		runNodeServer(ln)
	}()
	go func() {
		defer wg.Done()
		connectHeadServer(headAddress, nodeAddress)
	}()

	wg.Wait()
}

func connectHeadServer(headAddress string, nodeAddress string) {
	connHead, err := net.Dial("tcp", headAddress)
	if err != nil {
		fmt.Printf("통괄서버 연결 실패: %s\n", err)
		return
	}

	_, err = connHead.Write([]byte(nodeAddress))
	if err != nil {
		fmt.Printf("서버에게 ListenPort 전달 실패: %s\n", err)
		return
	}

	connectNodes(connHead, nodeAddress)

	fmt.Println("작동 준비 완료")

	//서버로부터 명령 받아
	//"전체" 이면 "전체" 출력 후 모든 노드에 "전체" 메시지 전달
	//"단일" 이면 "단일" 출력
	for {
		message, err := bufio.NewReader(connHead).ReadString('\n')
		if err != nil {
			fmt.Printf("서버로부터 입력 에러:%s\n", err)
			return
		}
		message = strings.TrimSpace(message)

		switch message {
		case "전체":
			fmt.Println("전체")
			//다른 노드로 전달
			for _, node := range nodes {
				node.Write([]byte(message + "\n"))
			}
		case "단일":
			fmt.Println("단일")
		default:
			fmt.Println("뭥미")
		}
	}
}

func connectNodes(connHead net.Conn, nodeAddress string) {
	//연결중인 노드 정보 받아오기
	nodeData := make([]byte, 1024)
	n, err := connHead.Read(nodeData)
	if err != nil {
		fmt.Printf("nodes 수신 오류: %s\n", err)
		return
	}

	var receivedMap map[string]map[string]string
	err = json.Unmarshal(nodeData[:n], &receivedMap)
	if err != nil {
		fmt.Printf("노드 데이터 역직렬화 실패:%s\n", err)
		return
	}

	//실행 중인 노드에 커넥트
	for _, nodeData := range receivedMap {
		connectNodeListener(nodeData["listenPort"], nodeAddress)
	}
}

func connectNodeListener(nodeListenPort string, thisListenPort string) {
	connNode, err := net.Dial("tcp", nodeListenPort)
	if err != nil {
		fmt.Printf("노드(%s) 연결 실패:%s\n", nodeListenPort, err)
		return
	}

	//승인해준 노드에게 자신의 ListenPort 전송
	_, err = connNode.Write([]byte(thisListenPort))
	if err != nil {
		fmt.Printf("노드(%s)에게 ListenPort전달 실패:%s\n", nodeListenPort, err)
		return
	}

	fmt.Printf("노드(%s)와 연결\n", nodeListenPort)
	nodes[nodeListenPort] = connNode

	go handleNode(connNode, nodeListenPort)
}

func runNodeServer(ln net.Listener) {
	fmt.Println("Starting...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("노드 연결 실패:", err)
			continue
		}

		//승인한 노드의 ListenPort 받아오기
		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Printf("노드 ListenPort 입력 오류:%s\n", err)
			return
		}
		listenPort := string(buffer[:bytesRead])

		fmt.Printf("노드(%s)와 연결\n", listenPort)
		nodes[listenPort] = conn

		go handleNode(conn, listenPort)
	}
}

func handleNode(conn net.Conn, nodeAddress string) {
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		message = strings.TrimSpace(message)
		if err != nil {
			fmt.Printf("노드(%s)와 연결끊김\n", nodeAddress)
			return
		}

		fmt.Printf("노드(%s)로부터의 message:%s\n", nodeAddress, message)
	}
}
