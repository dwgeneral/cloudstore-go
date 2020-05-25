package db

import (
	mydb "cloudstore-go/db/mysql"
	"database/sql"
	"fmt"
)

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

// OnFileUploadFinished 文件上传完成持久化到数据库
func OnFileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	statement, err := mydb.DBConn().Prepare("insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values(?,?,?,?,1)")
	if err != nil {
		fmt.Printf("Failed to prepare statement, err: %s\n", err.Error())
		return false
	}
	defer statement.Close()
	ret, err := statement.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Printf("failed to statement exec, err: %s\n", err.Error())
		return false
	}
	if returnFlag, err := ret.RowsAffected(); nil == err {
		if returnFlag <= 0 {
			fmt.Printf("File with hash: %s has been uploaded before\n", filehash)
		}
		return true
	}
	return false
}

// GetFileMeta 从数据库读取文件元信息
func GetFileMeta(filehash string) (*TableFile, error) {
	statement, err := mydb.DBConn().Prepare(
		"select file_sha1,file_addr,file_name,file_size from tbl_file where filesha1=? and state=1 limit 1")
	if err != nil {
		fmt.Printf("failed to query, err: %s", err.Error())
		return nil, err
	}
	defer statement.Close()
	tfile := TableFile{}
	err = statement.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileName, &tfile.FileSize)
	if err != nil {
		fmt.Printf("failed to queryRow, err: %s", err.Error())
		return nil, err
	}
	return &tfile, nil
}
