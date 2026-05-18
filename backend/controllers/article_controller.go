package controllers

import (
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 此文件用于实现创建文章的全部功能
func CreateArticle(ctx *gin.Context) {
	//创建一个Article结构体实例，并初始化其字段
	var article models.Article
	//将请求体中的JSON数据绑定到article结构体实例中
	if err := ctx.ShouldBindJSON(&article); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	//使用全局数据库连接对象global.Db调用AutoMigrate方法自动迁移article模型，确保数据库中有对应的表结构
	if err := global.Db.AutoMigrate(&article); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//调用Create方法将article结构体实例保存到数据库中，如果保存失败，返回500错误和错误信息
	if err := global.Db.Create(&article).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//如果保存成功，返回200状态码和保存的文章数据
	ctx.JSON(http.StatusOK, article)
}

// 获取文章列表函数
func GetArticles(ctx *gin.Context) {
	//拿到全部文章，需要用到一个Article结构体切片来存储查询结果
	var articles []models.Article
	//使用全局数据库连接对象global.Db调用Find方法查询所有的文章数据，并将结果保存到articles切片中，如果查询失败，返回500错误和错误信息
	if err := global.Db.Find(&articles).Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	//如果查询成功，返回200状态码和查询到的文章数据
	ctx.JSON(http.StatusOK, articles)
}
func GetArticleByID(ctx *gin.Context) {
	//从URL参数中获取文章ID
	//ctx.Param("id")获取URL参数中的id值，并将其转换为整数类型，如果转换失败，返回400错误和错误信息
	id := ctx.Param("id")
	var article models.Article
	if err := global.Db.Where("id=?", id).First(&article).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": "Article not found",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

	}
	ctx.JSON(http.StatusOK, article)
}
