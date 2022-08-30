package main

import (
	"net/http"
)

func main() {
	// http.Handle("/static/",
	// 	http.StripPrefix("/static/", http.FileServer(http.Dir("../../static"))))

	// http.HandleFunc("/file/upload", handler.UploadHandler)
	// http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	// http.HandleFunc("/file/meta", handler.GetFileMetaHandler)
	// http.HandleFunc("/file/download", handler.DownloadHandler)
	// http.HandleFunc("/file/update", handler.UpdateFileMetaHandler)
	// http.HandleFunc("/file/delete", handler.FileDeleteHandle)

	// http.HandleFunc("/user/signup", handler.UserSignUpHandler)
	// http.HandleFunc("/user/signin", handler.UserSignInHandler)
	// http.HandleFunc("/user/info", handler.RequestInterceptor(handler.UserInfoHandler))
	// http.HandleFunc("/file/mpupload/init", handler.RequestInterceptor(handler.MutilPartUploaInitdHandler))
	// http.HandleFunc("/file/mpupload/uppart", handler.RequestInterceptor(handler.UploadPartHandler))
	// http.HandleFunc("/file/mpupload/complete", handler.RequestInterceptor(handler.CompleteUplaodHandler))
	
	err := http.ListenAndServe(":8888", nil)

	if err != nil {
		panic(err)
	}

}
