package db

import (
	"fmt"
	mydb "go_code/project14/file-store-server/db/mysql"
)

type UserFile struct {
	UserName string
	FileSha1 string
	FileSize int64
	FileName string
	UploadAt string
}

func OnUserFileUpdateFinishied(username, filehash, filename string, filesize int64) bool {
	sqlStr := "insert into `tbl_user_file`(user_name,file_sha1,file_name,file_size)" +
		" values(?,?,?,?)"
	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		fmt.Println("OnUserFileUpdateFinishied Prepare failed:", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, filehash, filename, filesize)
	if err != nil {
		fmt.Println("OnUserFileUpdateFinishied Exec failed:", err)
		return false
	}
	if rowAffected, err := ret.RowsAffected(); err != nil || rowAffected == 0 {

		fmt.Println("OnUserFileUpdateFinishied RowsAffected failed:", err)
		return false
	}
	return true
}

func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	sqlStr := "select file_sha1,file_name,file_size,update_at from `tbl_user_file` where user_name = ? limit ?"
	st, err := mydb.DBConn().Prepare(sqlStr)
	defer st.Close()

	if err != nil {
		fmt.Println("QueryUserFileMetas Prepare failed:", err)
		return nil, err
	}

	rows, err := st.Query(username, limit)

	if err != nil {
		fmt.Println("QueryUserFileMetas Query failed:", err)
		return nil, err
	}

	res := make([]UserFile, 0)

	for rows.Next() {
		u := UserFile{}
		err := rows.Scan(&u.FileSha1, &u.FileName, &u.FileSize,&u.UploadAt)
		if err != nil {
			fmt.Println(" rows.Scan :", err.Error())
			break
		}
		res = append(res, u)
	}
	return res, nil
}



func OnUserFileCancel(username, filehash, filename string) bool {
	sqlStr := "insert into `tbl_user_file`(user_name,file_sha1,file_name,status)" +
		" values(?,?,?,?)"
	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		fmt.Println("OnUserFileUpdateFinishied Prepare failed:", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(username, filehash, filename,1)
	if err != nil {
		fmt.Println("OnUserFileCancel Exec failed:", err)
		return false
	}
	if rowAffected, err := ret.RowsAffected(); err != nil || rowAffected == 0 {
		fmt.Println("OnUserFileCancel RowsAffected failed:", err)
		return false
	}
	return true
}