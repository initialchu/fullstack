package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	//一个简单的路由组
	auth := r.Group("/api/auth")
	{
		auth.POST("/login", controllers.Login)

		auth.POST("/register", controllers.Register)
	}

	api := r.Group("/api")
	//不需要jwt认证的路由
	api.GET("/exchangerate", controllers.GetExchangeRate)
	//需要jwt认证的路由
	api.Use(middlewares.AuthMiddleware())
	{
		api.POST("/exchangerate", controllers.CreateExchangeRate)
	}
	return r

}
