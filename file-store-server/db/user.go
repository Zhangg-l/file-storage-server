package db

import (
	"fmt"
	mydb "go_code/project14/file-store-server/db/mysql"
)

func UserSignUp(username, password string, phone string) bool {
	dbStr := "insert ignore into tbl_user(`user_name`,`user_pwd`,`phone`) values(?,?,?)"

	stmt, err := mydb.DBConn().Prepare(dbStr)
	if err != nil {
		fmt.Println(" mydb.DBConn().Prepare faild:", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, password, phone)
	if err != nil {
		fmt.Println("stmt.Exec:", err)
		return false
	}

	if row, err := ret.RowsAffected(); err == nil && row > 0 {
		return true
	}

	return false
}

func UserSignIn(username, password string) bool {
	dbStr := "select `user_pwd` from tbl_user where user_name = ?"

	stmt, err := mydb.DBConn().Prepare(dbStr)
	if err != nil {
		fmt.Println(" mydb.DBConn().Prepare faild:", err)
		return false
	}
	defer stmt.Close()
	row := stmt.QueryRow(username)
	var pwd string
	if err := row.Scan(&pwd); err != nil {
		fmt.Println(" row.Scan:", err)
		return false
	}
	if pwd == password {
		return true
	}
	return false
}

func UserUpdateToken(username, token string) bool {
	dbStr := "replace into tbl_user_token(user_name,user_token) values(?,?)"

	stmt, err := mydb.DBConn().Prepare(dbStr)
	if err != nil {
		fmt.Println(" mydb.DBConn().Prepare faild:", err)
		return false
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, token)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}

func GetUserToken(username string) string {
	dbStr := "select user_token from `tbl_user_token` where user_name = ?"
	stmt, err := mydb.DBConn().Prepare(dbStr)
	if err != nil {
		fmt.Println(" mydb.DBConn().Prepare faild:", err)
		return ""
	}
	defer stmt.Close()

	row := stmt.QueryRow(username)
	var res string
	if err := row.Scan(&res); err != nil {
		fmt.Println("row.Scan:", err)
		return ""
	}
	return res
}

type UserInfo struct {
	Username string
	Email    string
	Phone    string
	SignupAt string
	Status   int
}

func GetUserInfo(username string) (UserInfo, error) {
	dbStr := "select user_name,signup_at from `tbl_user` where user_name = ?  limit 1"
	stmt, err := mydb.DBConn().Prepare(dbStr)
	userInfo := UserInfo{}
	if err != nil {
		fmt.Println("Prepare err:", err)
		return userInfo, err
	}
	row := stmt.QueryRow(username)
	if err = row.Scan(&userInfo.Username, &userInfo.SignupAt); err != nil {
		fmt.Println("Scan err:", err)
		return userInfo, err
	}

	return userInfo, nil
}
