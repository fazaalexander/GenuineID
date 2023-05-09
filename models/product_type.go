package models

import "gorm.io/gorm"

type Product_Type struct {
	gorm.Model
	Name     string    `json:"name" form:"name"`
	Products []Product `json:"products" form:"-" gorm:"foreignkey:Product_type_id"`
}
