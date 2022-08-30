package route

import (
	"go_code/project14/file-store-server/handler"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.Static("/static/", "./static")
	r.GET("/user/signup", handler.UserSignUpHandler)
	r.GET("/user/signin", handler.UserSignInHandler)
	r.POST("/user/signup", handler.DoUserSignUpHandler)
	r.POST("/user/signin", handler.DoUserSignInHandler)

	r.Use(handler.RequestInterceptor())
	r.GET("/file/upload", handler.UploadHandler)
	r.GET("/file/upload/suc", handler.UploadSucHandler)
	r.GET("/file/meta", handler.GetFileMetaHandler)
	r.GET("/file/download", handler.DownloadHandler)
	r.POST("/file/update", handler.UpdateFileMetaHandler)
	r.DELETE("/file/delete", handler.FileDeleteHandle)

	// r.GET("/user/info", handler.UserInfoHandler)
	r.GET("/file/mpupload/init", (handler.MutilPartUploaInitdHandler))
	r.POST("/file/mpupload/uppart", (handler.UploadPartHandler))
	r.GET("/file/mpupload/complete", (handler.CompleteUplaodHandler))

	return r
}
