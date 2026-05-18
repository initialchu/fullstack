package controllers

import (
	"encoding/json"
	"errors"
	"exchangeapp/global"
	"exchangeapp/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var cacheKey = "articles" //用于实现旁路缓存的键名

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
	cachedData, err := global.RedisDB.Get(cacheKey).Result()
	if err == redis.Nil {

		//拿到全部文章，需要用到一个Article结构体切片来存储查询结果
		var articles []models.Article
		//使用全局数据库连接对象global.Db调用Find方法查询所有的文章数据，并将结果保存到articles切片中，如果查询失败，返回500错误和错误信息
		if err := global.Db.Find(&articles).Error; err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
		//这里是缓存未命中，将查询到的文章数据转换为JSON格式，并将其存储在Redis中，设置过期时间为10分钟，如果存储失败，返回500错误和错误信息
		articleJSON, err := json.Marshal(articles)
		if err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		if err := global.RedisDB.Set(cacheKey, articleJSON, 10*time.Minute).Err(); err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(200, articles)
	} else if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	} else {
		//缓存命中
		var articles []models.Article
		//将从Redis中获取到的JSON数据反序列化为Article结构体切片，如果反序列化失败，返回500错误和错误信息
		if err := json.Unmarshal([]byte(cachedData), &articles); err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		//删除缓存中的数据，确保下一次请求能够获取到最新的文章数据
		if err := global.RedisDB.Del(cacheKey).Err(); err != nil {
			ctx.JSON(500, gin.H{
				"error": err.Error(),
			})
		}
		ctx.JSON(http.StatusOK, articles)
	}

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
