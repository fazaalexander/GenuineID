package models

import (
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
	Cust_id uint `json:"cust_id" form:"cust_id"`
	// Payment_id int `json:"payment_id" form:"payment_id"`
	Order_date          time.Time            `json:"order_date" form:"order_date"`
	Total_price         float64              `json:"total_price" form:"total_price"`
	Payment_status      bool                 `json:"payment_status" form:"payment_status"`
	Transaction_Details []Transaction_Detail `json:"transaction_details" form:"-" gorm:"foreignkey:Transaction_id"`
	Customer            Customer             `gorm:"foreignkey:Cust_id"`
	// Payment    Payment  `gorm:"foreignkey:Payment_id"`
}
