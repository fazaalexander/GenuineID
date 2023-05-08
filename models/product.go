package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Seller_id       uint             `json:"seller_id" form:"seller_id"`
	Product_type_id uint             `json:"product_type_id" form:"product_type_id"`
	Name            string           `json:"name" form:"name"`
	Description     string           `json:"description" form:"description"`
	Price           float64          `json:"price" form:"price"`
	Is_Verified     bool             `json:"is_verified" form:"is_verified"`
	Reviews         []Product_Review `json:"reviews" form:"-" gorm:"foreignkey:Product_id"`
	Seller          Seller           `gorm:"foreignKey:Seller_id"`
	Product_Type    Product_Type     `gorm:"foreignKey:Product_type_id"`
}

type ProductResponse struct {
	gorm.Model
	ID          uint    `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	IsVerified  bool    `json:"is_verified"`
}
