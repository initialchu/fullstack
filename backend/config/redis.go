package config

import (
	"exchangeapp/global"
	"log"

	"github.com/go-redis/redis"
)

func InitRedis() {
	//创建一个Redis客户端实例，连接到本地的Redis服务器
	RedisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // Redis服务器地址和端口
		DB:       0,                // 使用默认的数据库
		Password: "",               // 如果Redis服务器设置了密码，在此处提供
	})
	_, err := RedisClient.Ping().Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	global.RedisDB = RedisClient

}
