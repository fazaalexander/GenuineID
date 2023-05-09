package models

import "gorm.io/gorm"

type Transaction_Detail struct {
	gorm.Model
	Product_id     uint        `json:"product_id" form:"product_id"`
	Transaction_id uint        `json:"transaction_id" form:"transaction_id"`
	Seller_id      uint        `json:"seller_id" form:"seller_id"`
	Qty            int         `json:"qty" form:"qty"`
	Price          float64     `json:"price" form:"price"`
	Product        Product     `gorm:"foreignkey:Product_id"`
	Transaction    Transaction `gorm:"foreignkey:Transaction_id"`
}

type Transaction_Detail_Request struct {
	Product_id uint `json:"product_id" form:"product_id"`
	Qty        int  `json:"qty" form:"qty"`
}
