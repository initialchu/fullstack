package models

import "time"

//此文件用于定义汇率数据的结构体模型
type ExchangeRate struct {
	//ID字段，使用gorm的primarykey标签标记为主键，并且在JSON序列化时使用"id"作为字段名
	ID uint `gorm:"primarykey" json:"id"`
	//FromCurrency字段，表示汇率的来源货币，在JSON序列化时使用"from_currency"作为字段名
	FromCurrency string `json:"fromCurrency" binding:"required"` //必填字段
	//ToCurrency字段，表示汇率的目标货币，在JSON序列化时使用"to_currency"作为字段名
	ToCurrency string `json:"toCurrency" binding:"required"` //必填字段
	//Rate字段，表示汇率的数值，在JSON序列化时使用"rate"作为字段名
	Rate float64   `json:"rate" binding:"required"` //必填字段
	Date time.Time `json:"date"`
}
