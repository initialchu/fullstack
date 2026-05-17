package controllers

import (
	"exchangeapp/models"
	"exchangeapp/utils"

	"github.com/gin-gonic/gin"
)

func Register(ctx *gin.Context) {
	var user models.User
	//将请求体中的JSON数据绑定到user结构体实例中
	//使用gin的ShouldBindJSON方法，如果绑定失败，返回400错误和错误信息
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}
	//调用utils包中的HashPassword函数对用户输入的密码进行哈希处理，如果哈希处理失败，返回500错误和错误信息
	hashedPwd, err := utils.HashPassword(user.Password)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	//将哈希后的密码赋值回user结构体的Password字段
	user.Password = hashedPwd
}
