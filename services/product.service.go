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
