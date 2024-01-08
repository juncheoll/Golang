package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"node/database"
)

func FileListHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	files, err := database.GetFilesFromDB(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 파일 목록을 JSON 형식으로 반환
	jsonResponse, err := json.Marshal(files)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}
