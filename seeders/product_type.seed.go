package seeders

import (
	"github.com/fazaalexander/GenuineID/config"
	"github.com/fazaalexander/GenuineID/models"
)

func ProductTypeSeed() error {
	productTypes := []models.Product_Type{
		{Name: "Electronics"},
		{Name: "Fashion"},
		{Name: "Home & Kitchen"},
		{Name: "Jewelry"},
	}

	for _, pt := range productTypes {
		var productType models.Product_Type
		result := config.DB.Where("name = ?", pt.Name).First(&productType)
		if result.RowsAffected == 0 {
			if err := config.DB.Create(&pt).Error; err != nil {
				return err
			}
		}
	}

	return nil
}
