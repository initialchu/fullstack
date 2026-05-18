package main

import (
	//注意“exchangeapp”是go.mod文件中定义的模块名
	"context"
	"exchangeapp/config" //引入config包
	"exchangeapp/router"
	"log"

	//引入gin: go get github.com/gin-gonic/gin
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	//加载配置
	config.InitConfig()
	r := router.SetupRouter() //接受router包中SetupRouter函数返回的gin.Engine实例
	srv := &http.Server{
		Addr:    config.AppConfig.App.Port, //从配置文件中获取端口号
		Handler: r,
	}
	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no params) by default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器 ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Println("Server Shutdown:", err)
	}
	log.Println("服务器已关闭")
}
