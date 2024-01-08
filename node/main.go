package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"node/database"
	"node/handlers"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		fmt.Println("arguments 부족")
		return
	}

	portStr := args[1]
	dbname := args[2]
	port, err := strconv.Atoi(portStr)
	if err != nil {
		fmt.Println("포트 번호 숫자로 입력하시오")
		return
	}

	if dbname != "node1" && dbname != "node2" && dbname != "node3" && dbname != "node4" {
		fmt.Println("데이터베이스 이름 다시 쓰기!")
		return
	}
	fmt.Printf("포트 번호: %d, 데이터베이스 이름: %s\n", port, dbname)

	db := database.InitDB(dbname)
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handlers.IndexHandler(w, r, db)
	})
	http.HandleFunc("/filelist", func(w http.ResponseWriter, r *http.Request) {
		handlers.FileListHandler(w, r, db)
	})
	http.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		handlers.UploadHandler(w, r, db, dbname)
	})
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		handlers.DownloadHandler(w, r, db)
	})

	fmt.Printf("Server is running on port %d\n", port)
	http.ListenAndServe(":"+portStr, nil)
}
