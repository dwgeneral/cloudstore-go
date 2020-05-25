package main

import (
	"cloudstore-go/handler"
	"fmt"
	"net/http"
)

func main() {
	// 处理静态资源映射
	http.Handle("/static/",http.StripPrefix("/static",http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/file/upload", handler.UploadHandler)         // 创建上传路由
	http.HandleFunc("/file/upload/success", handler.UploadSuccess) // 上传成功路由
	http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	http.HandleFunc("/file/download", handler.DownloadHandler)
	http.HandleFunc("/file/update", handler.UpdateFileMetaHandler)
	http.HandleFunc("/file/delete", handler.DeleteFileHandler)
	http.HandleFunc("/user/signup", handler.SignupHandler)
	http.HandleFunc("/user/signin", handler.SignInHandler)
	fmt.Println("Server started on http://localhost:8080/user/signin")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %s", err.Error())
	}
}
