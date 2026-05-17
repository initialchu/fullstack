package main

import (
	//注意“exchangeapp”是go.mod文件中定义的模块名
	"exchangeapp/config" //引入config包
	"fmt"

	//引入gin: go get github.com/gin-gonic/gin
	"github.com/gin-gonic/gin"
)

func main() {
	//加载配置
	config.InitConfig()
	fmt.Println(config.AppConfig.App.Port)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run(config.AppConfig.App.Port) //监听配置文件中定义的端口
}
