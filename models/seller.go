package models

import (
	"gorm.io/gorm"
)

type Seller struct {
	gorm.Model
	User_id    uint      `json:"user_id" form:"user_id"`
	Store_name string    `json:"store_name" form:"store_name"`
	User       User      `gorm:"foreignKey:User_id"`
	Products   []Product `json:"products" form:"-" gorm:"foreignkey:Seller_id"`
}
