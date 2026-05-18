package controllers

import (
	"exchangeapp/global"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
)

// LikeArticle 处理点赞相关的请求
func LikeArticle(ctx *gin.Context) {
	// 从URL参数中获取文章ID
	articleID := ctx.Param("id")
	// 在这里可以添加逻辑来处理点赞，例如更新数据库中的点赞数

	likeKey := "article:" + articleID + ":likes"

	if err := global.RedisDB.Incr(likeKey).Err(); err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "点赞成功",
	})

}

// GetLikes 获取文章的点赞数
func GetLikes(ctx *gin.Context) {
	articleID := ctx.Param("id")
	likeKey := "article:" + articleID + ":likes"
	likes, err := global.RedisDB.Get(likeKey).Result()
	//如果Redis中没有该键，Get方法会返回一个redis.Nil错误，这时我们可以将点赞数设置为0；如果发生其他错误，则返回500错误和错误信息；如果成功获取点赞数，则返回200状态码和点赞数。
	if err == redis.Nil {
		likes = "0"
	} else if err != nil {
		ctx.JSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"likes": likes,
	})
}
