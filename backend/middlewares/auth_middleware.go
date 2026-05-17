package middlewares

import (
	"exchangeapp/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 此文件用于实现JWT认证中间件的全部功能
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//在这里实现JWT认证逻辑，例如从请求头中获取JWT令牌，验证令牌的有效性等
		token := ctx.GetHeader("Authorization")
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is missing",
			})
			//终止请求处理，返回401错误
			ctx.Abort()
			return
		}
		username, err := utils.ParseJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			//如果c.Next()之前调用了ctx.Abort()，则后续的处理函数将不会被执行，直接返回响应给客户端
			ctx.Abort()
			return
		}
		//将解析出的用户名存储在请求上下文中，以便后续处理函数使用
		//人话：把用户名放在上下文里，让其他函数夸中间件访问
		//通过ctx.Get("username")可以在后续的处理函数中获取到这个用户名
		ctx.Set("username", username)
		//继续处理请求，调用下一个处理函数
		ctx.Next()
	}
}
