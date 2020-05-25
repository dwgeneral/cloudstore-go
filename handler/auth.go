package handler

import (
	"net/http"
)

// HTTPInterceptor : http请求拦截器
func HTTPInterceptor(handleFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()
			username := r.Form.Get("username")
			token := r.Form.Get("token")
			//验证登录token是否有效
			if len(username) < 3 || !IsTokenValid(token) {
				http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
				return
			}
			handleFunc(w, r)
		})
}
