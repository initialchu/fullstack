package main

import (
	//注意“exchangeapp”是go.mod文件中定义的模块名
	"exchangeapp/config" //引入config包
	"exchangeapp/router"
	//引入gin: go get github.com/gin-gonic/gin
)

func main() {
	//加载配置
	config.InitConfig()
	r := router.SetupRouter() //接受router包中SetupRouter函数返回的gin.Engine实例
	port := config.AppConfig.App.Port
	if port == "" {
		port = ":3000" //如果配置文件中没有定义端口，使用默认端口3000
	}
	r.Run(config.AppConfig.App.Port) //监听配置文件中定义的端口
}
