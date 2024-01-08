package handlers

import (
	"database/sql"
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB, dbname string) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	r.ParseMultipartForm(10 << 20) // 최대 파일 크기: 10MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error retrieving the file")
		fmt.Println(err)
		return
	}
	defer file.Close()

	fmt.Println()

	// 업로드한 파일 저장
	uploadedFileName := handler.Filename
	filePath := "C:/Users/th6re8e/OneDrive - 계명대학교/Golang/node/uploads/" + dbname + "/" + uploadedFileName

	dst, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Error creating the file")
		fmt.Println(err)
		return
	}
	defer dst.Close()

	io.Copy(dst, file)

	// 파일 경로를 데이터베이스에 저장
	_, err = db.Exec("INSERT INTO files (filename, filepath) VALUES (?, ?);", uploadedFileName, filePath)
	if err != nil {
		fmt.Println("Error inserting into database")
		fmt.Println(err)
		return
	}

	fmt.Println("업로드 완료")
}
