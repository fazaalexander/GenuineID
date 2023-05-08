// Branch feature/product
package services

import (
	"net/http"
	"time"

	"github.com/fazaalexander/GenuineID/config"
	"github.com/fazaalexander/GenuineID/models"
	"github.com/labstack/echo/v4"
)

// Interface Product
type IProductService interface {
	CreateProduct(product models.Product) error
	GetProducts() (*[]models.ProductResponse, error)
	GetProductByID(productID string) (*models.ProductResponse, error)
	DeleteProduct(productID string) error
	UpdateProduct(productID string, name string, description string, price float64) error
	AuthenticateProduct(adminID uint, product_auth *models.Product_Auth) (*models.Product_Auth, error)
}

type ProductRepository struct {
	Func IProductService
}

var productRepository IProductService

func init() {
	ur := &ProductRepository{}
	ur.Func = ur

	productRepository = ur
}

func GetProductRepository() IProductService {
	return productRepository
}

func SetProductRepository(ur IProductService) {
	productRepository = ur
}

func (p *ProductRepository) CreateProduct(product models.Product) error {
	if err := config.DB.Save(&product).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Menambahkan data ke tabel product_auth
	productAuth := &models.Product_Auth{
		Admin_id:    nil,
		Product_id:  product.ID,
		Auth_date:   time.Now(),
		Auth_status: "Pending",
	}

	if err := config.DB.Create(&productAuth).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return nil
}

func (p *ProductRepository) GetProducts() (*[]models.ProductResponse, error) {
	var products []models.ProductResponse

	if err := config.DB.Model(&models.Product{}).Select("id, name, description, price, is_verified").Find(&products).Error; err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return &products, nil
}

func (p *ProductRepository) GetProductByID(productID string) (*models.ProductResponse, error) {
	var product models.ProductResponse

	if err := config.DB.Model(&models.Product{}).Select("id, name, description, price, is_verified").First(&product, productID).Error; err != nil {
		return &product, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return &product, nil
}

func (p *ProductRepository) DeleteProduct(productID string) error {
	var product models.Product

	if err := config.DB.First(&product, productID).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Delete product reviews
	if err := config.DB.Where("product_id = ?", productID).Delete(&[]models.Product_Review{}).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "Failed to delete product reviews",
		})
	}

	// Delete product
	if err := config.DB.Delete(&product, productID).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "Failed to delete product",
		})
	}

	return nil
}

func (p *ProductRepository) UpdateProduct(productID string, name string, description string, price float64) error {
	var product models.Product

	if err := config.DB.Model(&product).Where("id = ?", productID).Updates(models.Product{Name: name, Description: description, Price: price}).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Failed update user",
		})
	}

	return nil
}

func (p *ProductRepository) AuthenticateProduct(admin_id uint, product_auth *models.Product_Auth) (*models.Product_Auth, error) {
	var product models.Product
	var admin models.Admin

	if err := config.DB.Where("user_id = ?", admin_id).First(&admin).Error; err != nil {
		return product_auth, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Admin not found",
		})
	}

	if err := config.DB.Where("id = ?", &product_auth.Product_id).First(&product).Error; err != nil {
		return product_auth, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Product not found",
		})
	}

	if err := config.DB.Model(&models.Product_Auth{}).Where("product_id = ?", product_auth.Product_id).Updates(models.Product_Auth{Admin_id: &admin.ID, Auth_status: product_auth.Auth_status}).Error; err != nil {
		return product_auth, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Failed to update auth status",
			"error":   err.Error(),
			"id":      &admin.ID,
		})
	}

	if product_auth.Auth_status == "original" {
		if err := config.DB.Model(&models.Product{}).Where("id = ?", product_auth.Product_id).Update("is_verified", true).Error; err != nil {
			return product_auth, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
				"message": "Failed to update product verified status",
			})
		}
	}

	if err := config.DB.Preload("Product").Find(&product_auth).Error; err != nil {
		return product_auth, echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return product_auth, nil
}
