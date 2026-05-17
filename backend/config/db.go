package config

import (
	"exchangeapp/global"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() {
	dsn := AppConfig.Database.Dsn
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database:%v", err)

	}
	sqlDB, err := db.DB()

	sqlDB.SetMaxIdleConns(AppConfig.Database.MaxIdleConns) //设置连接池中空闲连接的最大数量
	sqlDB.SetMaxOpenConns(AppConfig.Database.MaxOpenConns) //设置连接池中打开连接的最大数量
	sqlDB.SetConnMaxLifetime(time.Hour)                    //设置连接可复用的最长时间，0表示没有限制,这里是一小时

	if err != nil {
		log.Fatalf("failed to get database instance:%v", err)
	}
	global.Db = db
}
