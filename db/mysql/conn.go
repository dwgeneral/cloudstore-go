package mysql

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql" // method delegate to database/sql
)

var db *sql.DB

func init() {
	db, _ = sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/cloudstore?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		fmt.Printf("Failed to connect DB, err: %s\n", err.Error())
	}
}

// DBConn 返回数据库连接对象
func DBConn() *sql.DB {
	return db
}
