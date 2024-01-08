package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	fileName := r.URL.Query().Get("file")

	if fileName == "" {
		http.Error(w, "File name not specified", http.StatusBadRequest)
		return
	}
	fmt.Println(fileName, "다운로드 요청")

	// 파일이름을 통해 DB에서 파일 경로 조회
	var filePath string
	err := db.QueryRow("SELECT filepath FROM files WHERE filename = ?", fileName).Scan(&filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if filePath == "" {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	fmt.Println(filePath)
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	io.Copy(w, file)
}
