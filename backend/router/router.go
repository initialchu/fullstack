package router

import (
	"exchangeapp/controllers"
	"exchangeapp/middlewares"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	// 配置 CORS 中间件
	r.Use(cors.New(cors.Config{
		// 指定哪些来源可以发送请求到服务器
		AllowOrigins: []string{"http://localhost:5173"},
		// 允许哪些 HTTP 方法
		AllowMethods: []string{"GET", "POST", "OPTIONS"},
		// 允许哪些 HTTP 头部
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
		//哪些请求头可以暴露给浏览器
		ExposeHeaders: []string{"Content-Length"},
		// 是否允许发送 Cookie,如果前端需要发送认证信息（如 JWT token），则需要设置为 true
		AllowCredentials: true,
		// 预检请求的有效期，单位为秒
		MaxAge: 12 * 60 * 60,
	}))

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
		api.POST("/articles", controllers.CreateArticle)
		api.GET("/articles", controllers.GetArticles)
		api.GET("/articles/:id", controllers.GetArticleByID)
		api.GET("/articles/:id/like", controllers.GetLikes)
		api.POST("/articles/:id/like", controllers.LikeArticle)
	}
	return r

}
