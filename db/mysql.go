package db

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("mysql", "root:password@tcp(127.0.0.1:3306)/test")
	if err != nil {
		log.Fatal(err)
	}
}

// GetMySQLConnection will return the mysql collection
func GetMySQLConnection() *sql.DB {
	return db
}
