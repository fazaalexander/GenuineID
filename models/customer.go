package models

import (
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	User_id      uint             `json:"user_id" form:"user_id"`
	FullName     string           `json:"full_name" form:"full_name"`
	Transactions []Transaction    `gorm:"foreignKey:Cust_id"`
	Reviews      []Product_Review `gorm:"foreignKey:Cust_id"`
}
