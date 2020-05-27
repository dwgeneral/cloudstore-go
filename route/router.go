package route

import (
	"cloudstore-go/handler"

	"github.com/gin-gonic/gin"
)

// Router 路由
func Router() *gin.Engine {
	router := gin.Default()
	router.Static("/static/", "./static")

	router.GET("user/signup", handler.SignupHandler)
	router.POST("user/signup", handler.DoSignupHandler)

	router.GET("user/signin", handler.SignInHandler)
	router.POST("user/signin", handler.DosignInHandler)

	// 此行之后所有的 router 都会先走拦截器，进行token校验
	router.Use(handler.HTTPInterceptor())
}
