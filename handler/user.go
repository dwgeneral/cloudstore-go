package handler

import (
	"cloudstore-go/db"
	"cloudstore-go/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const salt = "less is more !"

// SignupHandler  用户注册处理
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 3 {
		w.Write([]byte("invalid parameters"))
		return
	}
	encodePasswd := util.Sha1([]byte(passwd + salt))
	result := db.UserSignUp(username, encodePasswd)
	if result {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("failed"))
	}
}

// SignInHandler 登陆处理
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	encodePasswd := util.Sha1([]byte(passwd + salt))

	pwdChecked := db.UserSignIn(username, encodePasswd)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}
	token := GenToken(username)
	res := db.UpdateToken(username, token)
	if !res {
		w.Write([]byte("failed"))
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
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
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
