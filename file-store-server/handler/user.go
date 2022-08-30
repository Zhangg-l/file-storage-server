package handler

import (
	"fmt"
	"go_code/project14/file-store-server/db"
	"go_code/project14/file-store-server/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	hash_salt = "@#*679"
)

func UserSignUpHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "./static/view/signup.html")
}

func DoUserSignUpHandler(c *gin.Context) {

	u_pwd := c.Request.FormValue("password")
	username := c.Request.FormValue("username")

	if len(u_pwd) < 3 || len(username) < 3 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "invalid paramer length",
			"code": "-1",
		})
		return
	}

	u_pwd = util.Sha1([]byte(u_pwd + hash_salt))
	phone := util.GetRandPhone()
	suc := db.UserSignUp(username, u_pwd, phone)
	if suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "SUCCESS",
			"code": "0",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg":  "FAILED",
		"code": "-1",
	})
	return
}
func UserSignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

func DoUserSignInHandler(c *gin.Context) {

	u_pwd := c.Request.FormValue("password")
	username := c.Request.FormValue("username")

	u_pwd = util.Sha1([]byte(u_pwd + hash_salt))
	ckeckSignIn := db.UserSignIn(username, u_pwd)
	if !ckeckSignIn {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "FAILED",
			"code": "-1",
		})
		return
	}
	// create token
	u_token := GenerationTakon(username)

	// save table
	suc := db.UserUpdateToken(username, u_token)
	if !suc {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "FAILED",
			"code": "-1",
		})
		return
	}
	// 登录成功后 重定向到首页
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Token:    u_token,
			Username: username,
			Location: "/static/view/home.html",
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

func GenerationTakon(username string) string {
	// token : username + timestampt + tokensalt + timestampt[:8]
	timestampt := fmt.Sprintln(time.Now().Unix())
	// len is 32
	_token := util.MD5([]byte(username + timestampt + "_tokensalt"))

	return _token + timestampt[:8]
}

func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	// u_token := r.Form.Get("token")

	// if !IsValidToken(u_token, username) {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }

	userInfo, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: userInfo,
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp.JSONBytes())
}

func IsValidToken(token string, username string) bool {
	if len(token) != 40 {
		return false
	}
	newTime := fmt.Sprintf("%d", time.Now().Unix())
	newTime = newTime[:8]
	oldTime := token[32:]

	otime, err := strconv.Atoi(oldTime)
	if err != nil {
		return false
	}
	ntime, err := strconv.Atoi(newTime)
	if err != nil {
		return false
	}

	// time expire
	if diff := ntime - otime; diff*100 > 100*60 {
		return false
	}
	// query
	u_token := db.GetUserToken(username)
	if token != u_token {
		return false
	}
	return true
}
