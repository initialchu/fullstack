package controllers

import (
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

//此文件用于实现创建汇率数据的全部功能

func CreateExchangeRate(ctx *gin.Context) {
	//创建一个ExchangeRate结构体实例，并初始化其字段
	var exchangeRate models.ExchangeRate
	//将请求体中的JSON数据绑定到exchangeRate结构体实例中
	if err := ctx.ShouldBindJSON(&exchangeRate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//将当前时间赋值给exchangeRate结构体的Date字段
	exchangeRate.Date = time.Now()
	//使用全局数据库连接对象global.Db调用AutoMigrate方法自动迁移exchangeRate模型，确保数据库中有对应的表结构
	if err := global.Db.AutoMigrate(&exchangeRate); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//调用Create方法将exchangeRate结构体实例保存到数据库中，如果保存失败，返回500错误和错误信息
	if err := global.Db.Create(&exchangeRate).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//如果保存成功，返回200状态码和保存的汇率数据
	ctx.JSON(http.StatusOK, exchangeRate)
}
