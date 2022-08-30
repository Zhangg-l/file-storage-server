package mysql

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	var err error

	db, err = sql.Open("mysql", "root:123456@tcp(127.0.0.1:3310)/filestore?charset=utf8")
	if err != nil {
		fmt.Println("failed to Open mysql:", err)
		return
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("failed to connection mysql:", err)
		os.Exit(1)
	}
}

func DBConn() *sql.DB {
	return db
}
