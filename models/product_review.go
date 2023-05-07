package models

import (
	"time"

	"gorm.io/gorm"
)

type Product_Review struct {
	gorm.Model
	Product_id  uint      `json:"product_id" form:"product_id"`
	Cust_id     uint      `json:"cust_id" form:"cust_id"`
	Review      string    `json:"review" form:"review"`
	Review_date time.Time `json:"review_date" form:"review_date"`
	Product     Product   `gorm:"foreignkey:Product_id"`
	Customer    Customer  `gorm:"foreignkey:Cust_id"`
}
