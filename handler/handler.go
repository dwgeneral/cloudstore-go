package handler

import (
	"cloudstore-go/meta"
	"cloudstore-go/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// UploadHandler 处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internal server error")
			return
		}
		io.WriteString(w, string(data))
	case "POST":
		// 接收文件流及存储到本地目录
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err: %s\n", err.Error())
			return
		}
		defer file.Close()

		// 存储文件元信息
		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			FilePath: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.FilePath)
		if err != nil {
			fmt.Printf("Failed to create file, err: %s\n", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file, err: %s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UploadFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/success", http.StatusFound)
	}
}

// UploadSuccess 上传成功
func UploadSuccess(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "upload success!")
}

// GetFileMetaHandler 根据 form hash 获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filehash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// DownloadHandler 下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fmeta := meta.GetFileMeta(fsha1)
	file, err := os.Open(fmeta.FilePath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fmeta.FileName+"\"")
	w.Write(data)
}

// UpdateFileMetaHandler 更新文件元信息
func UpdateFileMetaHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fileMeta := meta.GetFileMeta(filesha1)
	fileMeta.FileName = r.Form.Get("filename")
	meta.UploadFileMeta(fileMeta)
	data, err := json.Marshal(fileMeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

// DeleteFileHandler 删除文件及元信息数据
func DeleteFileHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	filesha1 := r.Form.Get("filehash")

	fmeta := meta.GetFileMeta(filesha1)
	os.Remove(fmeta.FilePath)

	meta.RemoveFileMeta(filesha1)
	w.WriteHeader(http.StatusOK)
}
