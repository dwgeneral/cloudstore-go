package handler

import (
	"cloudstore-go/db"
	"cloudstore-go/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const salt = "less is more !"

// DoSignupHandler 用户注册POST处理
func DoSignupHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	if len(username) < 3 || len(passwd) < 3 {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "invalid params",
			"code": -1,
		})
		return
	}
	encodePasswd := util.Sha1([]byte(passwd + salt))
	result := db.UserSignUp(username, encodePasswd)
	if result {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "success",
			"code": 0,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "sign up failed",
			"code": -2,
		})
	}
}

// SignupHandler  用户注册GET处理
func SignupHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signup.html")
}

// DosignInHandler 用户登陆处理POST
func DosignInHandler(c *gin.Context) {
	username := c.Request.FormValue("username")
	passwd := c.Request.FormValue("password")

	encodePasswd := util.Sha1([]byte(passwd + salt))

	pwdChecked := db.UserSignIn(username, encodePasswd)
	if !pwdChecked {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "login failed",
			"code": -1,
		})
		return
	}
	token := GenToken(username)
	res := db.UpdateToken(username, token)
	if !res {
		c.JSON(http.StatusOK, gin.H{
			"msg":  "update token failed when sign in",
			"code": -2,
		})
		return
	}
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	c.Data(http.StatusOK, "application/json", resp.JSONBytes())
}

// SignInHandler 登陆处理GET
func SignInHandler(c *gin.Context) {
	c.Redirect(http.StatusFound, "/static/view/signin.html")
}

// UserInfoHandler 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	user, err := db.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	res := util.RespMsg{
		Code: 0,
		Msg:  "ok",
		Data: user,
	}
	w.Write(res.JSONBytes())
}

// GenToken 生成token
func GenToken(username string) string {
	timestamp := fmt.Sprintf("%x", time.Now().Unix())
	prefix := util.MD5([]byte(username + timestamp + "_tokensalt"))
	return prefix + timestamp[:8]
}

// IsTokenValid 校验Token
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	return true
}
