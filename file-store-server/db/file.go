package db

import (
	"database/sql"
	"fmt"
	mydb "go_code/project14/file-store-server/db/mysql"
)

func OnUploadFileFinishied(fileSha1, fileName,
	fileAddre string, fileSize int64) bool {
	// perpare
	queryStr := "insert ignore into tbl_file " +
		"(`file_name`,`file_sha1`,`file_addr`,`file_size`,`status`) values(?,?,?,?,?)"
	stmt, err := mydb.DBConn().Prepare(queryStr)
	if err != nil {
		fmt.Println("failed to mysql.DBConn().Prepare:", err)
		return false
	}
	defer stmt.Close()
	ret, err := stmt.Exec(fileName, fileSha1, fileAddre, fileSize, 1)
	if err != nil {
		fmt.Println("failed to stmt.Exec:", err)
		return false
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf <= 0 {
			fmt.Printf("(Waring): File with hash:%s  has been uploaded before!\n", fileSha1)
			return false
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(fileSha1 string) (TableFile, error) {
	queryStr := "select file_name,file_addr,file_size from tbl_file where `file_sha1` = ? AND `status` = 1"
	stmt, err := mydb.DBConn().Prepare(queryStr)
	if err != nil {
		fmt.Println("failed to mysql.DBConn().Prepare:", err)
		return TableFile{}, err
	}
	defer stmt.Close()
	tf := TableFile{}
	err = stmt.QueryRow(fileSha1).Scan(&tf.FileName, &tf.FileAddr, &tf.FileSize)
	if err != nil {
		if err == sql.ErrNoRows {
			return TableFile{}, nil
		}
		fmt.Println("failed to QueryRow.Scan:", err)
		return TableFile{}, err
	}

	tf.FileHash = fileSha1
	return tf, nil
}

func UpdateFileLocation(filehash string, path string) error {
	sqlStr := "update `tbl_file` set `file_addr` = ? where `file_sha1` = ?"

	stmt, err := mydb.DBConn().Prepare(sqlStr)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	res, err := stmt.Exec(path,filehash)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	if rowAffect, err := res.RowsAffected(); err != nil || rowAffect <= 0 {
		fmt.Println(err.Error())
		return err
	}
	return nil
}
