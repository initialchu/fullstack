package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"exchangeapp/utils"
	"net/http"

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
			"error": "Failed to hash password" + err.Error(),
		})
		return
	}
	//将哈希后的密码赋值回user结构体的Password字段
	user.Password = hashedPwd

	//调用utils包中的GenerateJWT函数生成一个JWT令牌，如果生成失败，返回500错误和错误信息
	token, err := utils.GenerateJWT(user.Username)
	if err != nil {
		ctx.JSON(500, gin.H{
			"error": "Failed to generate token" + err.Error(),
		})
		return
	}
	//使用全局数据库连接对象global.Db调用AutoMigrate方法自动迁移user模型，确保数据库中有对应的表结构
	//人话：如果表不存在，自动创建一个用户表，已存在看有没有新字段需要添加
	if err := global.Db.AutoMigrate(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to migrate database",
		})
		return
	}
	//使用全局数据库连接对象global.Db调用Create方法将user结构体实例保存到数据库中，如果保存失败，返回500错误和错误信息
	//人话：创建一个用户
	if err := global.Db.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user" + err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"token": token,
	})
}
