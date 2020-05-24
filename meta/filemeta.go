package meta

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

// GetFileMeta 通过 Sha1 值获取文件的元信息对象
func GetFileMeta(fsha1 string) FileMeta {
	return fileMetas[fsha1]
}
