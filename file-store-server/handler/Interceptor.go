package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequestInterceptor() gin.HandlerFunc {

	return func(ctx *gin.Context) {
		username := ctx.Request.FormValue("username")
		u_token := ctx.Request.FormValue("token")
		if len(username) < 3 || !IsValidToken(u_token, username) {
			ctx.Abort()
			ctx.JSON(http.StatusOK, gin.H{
				"msg":  "ID Card invalid",
				"code": "-1",
			})
			return
		}
		ctx.Next()
	}
}
