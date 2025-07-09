package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("mysql", "root:kitchu@tcp(localhost:3307)/automation")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	err = DB.Ping()
	if err != nil {
		log.Fatal("Database ping error:", err)
	}
	log.Println("Database connected successfully")
}
