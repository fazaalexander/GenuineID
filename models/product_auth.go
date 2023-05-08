package models

import (
	"time"

	"gorm.io/gorm"
)

type Product_Auth struct {
	gorm.Model
	Admin_id    *uint     `json:"admin_id" form:"admin_id"`
	Product_id  uint      `json:"product_id" form:"product_id"`
	Auth_date   time.Time `json:"auth_date" form:"auth_date"`
	Auth_status string    `json:"auth_status" form:"auth_status"`
	Admin       Admin     `gorm:"foreignKey:Admin_id"`
	Product     Product   `gorm:"foreignKey:Product_id"`
}
