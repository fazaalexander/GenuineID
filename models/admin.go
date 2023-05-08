package models

import (
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	User_id uint `json:"user_id" form:"user_id"`

	User         User           `gorm:"foreignKey:User_id"`
	Product_Auth []Product_Auth `gorm:"foreignKey:Admin_id"`
}
