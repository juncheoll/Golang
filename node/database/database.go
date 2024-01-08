package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"fmt"
)

func InitDB(dbname string) *sql.DB {
	var db *sql.DB
	var err error

	db, err = sql.Open("mysql", "root:skaksdml59!@tcp(localhost:3306)/"+dbname)
	if err != nil {
		fmt.Println("Error opening database")
		return nil
	}

	err2 := db.Ping()
	if err2 != nil {
		fmt.Println("Connect Error")
		fmt.Println(err)
		return nil
	}

	// 파일 테이블 생성
	stmt, err := db.Prepare(`CREATE TABLE IF NOT EXISTS files (
		id INT AUTO_INCREMENT PRIMARY KEY,
		filename VARCHAR(255),
		filepath VARCHAR(255)
	)`)
	if err != nil {
		fmt.Println("Error creating table Query")
		return nil
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println("Error creating table")
		return nil
	}

	return db
}

func GetFilesFromDB(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT filename FROM files")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var files []string
	for rows.Next() {
		var filename string
		err := rows.Scan(&filename)
		if err != nil {
			return nil, err
		}
		files = append(files, filename)
	}

	return files, nil
}
