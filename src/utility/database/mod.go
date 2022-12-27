package database

import (
	"chimera/utility/configuration"
	"chimera/utility/logging"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func Initialize() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(127.0.0.1:3306)/phantom", configuration.GetConfiguration().Connection.DBLogin))

	if err != nil {
		logging.Fatal("Failed to contact database! (%s)", err.Error())
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	logging.Info("Database Service", "Connected to Database")
}

func Query(data string, args ...any) (*sql.Rows, error) {
	return db.Query(data, args...)
}
