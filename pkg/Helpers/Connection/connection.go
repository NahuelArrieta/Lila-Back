package connection

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	url := "root:root@tcp(localhost:3306)/api_videogame"
	var err error
	db, err = sql.Open("mysql", url)
	if err != nil {
		panic(err.Error())
	}
}

func Connect() (*sql.Tx, error) {
	txn, err := db.Begin()
	if err != nil {
		return nil, err
	}
	return txn, nil
}
