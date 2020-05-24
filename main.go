package main

import (
	"cloudstore-go/handler"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)         // 创建上传路由
	http.HandleFunc("/file/upload/success", handler.UploadSuccess) // 上传成功路由
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	fmt.Println("Server started on http://localhost:8080/file/upload")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %s", err.Error())
	}
}
