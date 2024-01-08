package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"text/template"
)

type NodeInfo struct {
	Port string
	Name string
}

var nodeInfoArr = []NodeInfo{
	{"9001", "node1"},
	{"9002", "node2"},
	{"9003", "node3"},
	{"9004", "node4"},
}

func main() {
	RunNodeServer()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("index.html"))

		err := tmpl.Execute(w, nodeInfoArr)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	fmt.Println("Client server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}

func RunNodeServer() {
	nodePath := "C:/Users/th6re8e/OneDrive - 계명대학교/Golang/node"

	for _, nodeInfo := range nodeInfoArr {
		cmd := exec.Command("cmd", "/C", "start", "cmd", "/K", "go", "run", nodePath, nodeInfo.Port, nodeInfo.Name)
		err := cmd.Run()
		if err != nil {
			fmt.Println("터미널 열기 실패:", err)
		}
	}

}
