package db

import (
	db "cloudstore-go/db/mysql"
	"fmt"
)

// UserSignUp 用户注册
func UserSignUp(username string, passwd string) bool {
	statement, err := db.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`, `user_pwd`) values (?,?)")
	if err != nil {
		fmt.Printf("Failed to insert, err: %s", err.Error())
		return false
	}
	defer statement.Close()

	ret, err := statement.Exec(username, passwd)
	if err != nil {
		fmt.Printf("Failed to insert, err: %s", err.Error())
		return false
	}
	if rowsAffected, err := ret.RowsAffected(); nil == err && rowsAffected > 0 {
		return true
	}
	return false
}

// UserSignIn 用户登陆鉴权
func UserSignIn(username string, encodePasswd string) bool {
	statement, err := db.DBConn().Prepare("select * from tbl_user where user_name = ? limit 1")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	rows, err := statement.Query(username)
	if err != nil || rows == nil {
		fmt.Println(err.Error())
		return false
	}

	pRows := db.ParseRows(rows)
	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encodePasswd {
		return true
	}
	return false
}

// UpdateToken 刷新用户Token
func UpdateToken(username string, token string) bool {
	statement, err := db.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`, `user_token`) values(?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer statement.Close()
	_, err = statement.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}
