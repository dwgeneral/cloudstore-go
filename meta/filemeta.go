package meta

import (
	db "cloudstore-go/db"
)

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	FilePath string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UploadFileMeta 新增/更新文件元信息
func UploadFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UploadFileMetaDB 更新文件元信息到数据库
func UploadFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.FilePath)
}

// GetFileMeta 通过 Sha1 值获取文件的元信息对象
func GetFileMeta(fsha1 string) FileMeta {
	return fileMetas[fsha1]
}

// GetFileMetaDB 通过 Sha1 值查询 tbl_file 表
func GetFileMetaDB(fsha1 string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(fsha1)
	if err != nil {
		return FileMeta{}, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		FilePath: tfile.FileAddr.String,
	}
	return fmeta, nil
}

// RemoveFileMeta 通过 Sha1 值删除文件的元信息对象
func RemoveFileMeta(fsha1 string) {
	delete(fileMetas, fsha1)
}
