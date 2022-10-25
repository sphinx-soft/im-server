package util

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDatabase() {
	database, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(127.0.0.1:3306)/phantom", GetDatabaseLogin()))
	db = database
	if err != nil {
		panic(err)
	}
	database.SetMaxOpenConns(100)
	database.SetMaxIdleConns(50)
	Log("Database", "Initialised database server")
}

func GetDatabaseHandle() *sql.DB {
	return db
}
