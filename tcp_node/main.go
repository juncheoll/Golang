package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
)

var nodes map[string]net.Conn

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("node 포트를 주세요")
		return
	}
	port := args[1]
	nodePort, err := strconv.Atoi(port)
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
		connectHeadServer(headAddress, nodePort)
	}()

	wg.Wait()
}

func connectHeadServer(headAddress string, nodePort int) {
	connHead, err := net.Dial("tcp", headAddress)
	if err != nil {
		fmt.Printf("통괄서버 연결 실패: %s\n", err)
		return
	}

	nodeData := make([]byte, 1024)
	n, err := connHead.Read(nodeData)
	if err != nil {
		fmt.Printf("nodes 수신 오류: %s\n", err)
		return
	}
	fmt.Print("nodes 크기:", n)

	var receiveNodeMap map[string]string
	err = json.Unmarshal(nodeData, &receiveNodeMap)
	if err != nil {
		fmt.Printf("노드 데이터 역직렬화 실패:%s\n", err)
		return
	}

	//실행 중인 노드에 커넥트
	nodes = make(map[string]net.Conn)
	for _, nodeAddress := range receiveNodeMap {
		connectNodeServer(nodeAddress)
	}

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

		switch message {
		case "전체":
			fmt.Println("전체")
			//다른 노드로 전달
			for _, node := range nodes {
				node.Write([]byte(message))
			}
		case "단일":
			fmt.Println("단일")
		default:
			fmt.Println("뭥미")
		}
	}
}

func connectNodeServer(nodeAddress string) {
	connNode, err := net.Dial("tcp", nodeAddress)
	if err != nil {
		fmt.Printf("노드(%s) 연결 실패:%s\n", nodeAddress, err)
		return
	}

	nodes[connNode.RemoteAddr().String()] = connNode
}

func runNodeServer(ln net.Listener) {
	fmt.Println("Starting...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("노드 연결 실패:", err)
			continue
		}
		fmt.Printf("다른 노드 연결: %s\n", conn.RemoteAddr().String())
		nodes[conn.RemoteAddr().String()] = conn
	}
}

func handleNode(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Printf("노드(%s)로부터 입력 오류:%s\n", conn.RemoteAddr().String(), err)
			return
		}

		fmt.Printf("노드(%s)로부터의 message:%s\n", conn.RemoteAddr().String(), message)
	}
}
