package database

import (
	"chimera/network"
	"chimera/utility/configuration"
	"chimera/utility/logging"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var database *sql.DB

func Initialize() {
	db, err := sql.Open("mysql", fmt.Sprintf("%s@tcp(127.0.0.1:3306)/chimera", configuration.GetConfiguration().Connection.DBLogin))

	if err != nil {
		logging.Fatal("Failed to contact database! (%s)", err.Error())
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(50)
	database = db
	logging.Info("Database Service", "Connected to Database")
}

func Query(data string, args ...any) (*sql.Rows, error) {
	return database.Query(data, args...)
}

func SetLastLoginDate(uin int) error {
	row, err := Query("UPDATE userdetails SET LastLogin= ? WHERE UIN= ?", time.Now().UnixNano(), uin)
	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return err
	}
	row.Close()
	return err
}

func GetAccountDataByEmail(email string) (network.Account, error) {

	var acc network.Account

	row, err := Query("SELECT * from accounts WHERE Mail= ?", email)

	if err != nil {
		logging.Error("Database/GetUserData", "Failed to get userdata: %s", err.Error())
		return acc, err
	}

	row.Next()
	row.Scan(&acc.UIN, &acc.DisplayName, &acc.Mail, &acc.Password)
	row.Close()

	return acc, err
}
