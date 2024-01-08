package handlers

import (
	"database/sql"
	"net/http"
	"node/database"
	"text/template"
)

func IndexHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	files, err := database.GetFilesFromDB(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl := template.Must(template.ParseFiles("C:/Users/th6re8e/OneDrive - 계명대학교/Golang/node/templates/index.html"))
	if err := tmpl.Execute(w, files); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
