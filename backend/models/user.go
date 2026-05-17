package models

import "gorm.io/gorm"

type User struct {
	/*
		GORM提供了一个预定义的结构体，名为gorm.Model，其中包含常用字段：

		// gorm.Model 的定义
		type   类型 Model struct {   模型结构{
	  ID        uint           `gorm:"primaryKey"`
	  ID        uint           `gorm:"primaryKey"`

	  CreatedAt time.Time   CreatedAt时间。时间
	  UpdatedAt time.Time   UpdatedAt时间。时间
	  DeletedAt gorm.DeletedAt `gorm:"index"`
	DeletedAt弄脏。DeletedAt“gorm:“index"”

	}

	*/
	gorm.Model
	Username string `gorm:"unique"` //用户名唯一
	Password string
}
