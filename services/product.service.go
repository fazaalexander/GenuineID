// Branch feature/product
package services

import (
	"log"
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
	SearchProductByName(name string) (*[]models.ProductResponse, error)
	SearchProductByID(product_id string) (*models.ProductResponse, error)
	SearchProductByType(product_type_id string) (*[]models.ProductResponse, error)
	ProductCheckout(cust_id uint, transaction_details []models.Transaction_Detail_Request) error
	CreateProductReview(cust_id uint, req models.Product_Review_Request) error
}

type ProductRepository struct {
	Func IProductService
}

var productRepository IProductService

func init() {
	pr := &ProductRepository{}
	pr.Func = pr

	productRepository = pr
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

func (p *ProductRepository) SearchProductByName(name string) (*[]models.ProductResponse, error) {
	var products []models.ProductResponse
	if err := config.DB.Model(&models.Product{}).Where("name LIKE ?", "%"+name+"%").Where("is_verified = true").Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Product not found",
		})
	}

	return &products, nil
}
func (p *ProductRepository) SearchProductByID(product_id string) (*models.ProductResponse, error) {
	var product *models.ProductResponse
	if err := config.DB.Model(&models.Product{}).Where("id = ?", product_id).Where("is_verified = true").First(&product).Error; err != nil {
		return nil, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Product not found",
		})
	}

	return product, nil
}

func (p *ProductRepository) SearchProductByType(product_type_id string) (*[]models.ProductResponse, error) {
	var products []models.ProductResponse
	if err := config.DB.Model(&models.Product{}).Where("product_type_id = ?", product_type_id).Where("is_verified = true").Find(&products).Error; err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Product not found",
		})
	}

	return &products, nil
}

func (p *ProductRepository) ProductCheckout(user_id uint, req []models.Transaction_Detail_Request) error {
	var totalPrice float64
	for _, td := range req {
		var product models.Product
		if err := config.DB.Model(&models.Product{}).Where("id = ?", td.Product_id).Where("is_verified = true").First(&product).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}
		totalPrice += float64(td.Qty) * product.Price
	}

	var customer models.Customer
	if err := config.DB.Model(&models.Customer{}).Where("user_id = ?", user_id).First(&customer).Error; err != nil {
		return err
	}

	transaction := models.Transaction{
		Cust_id:        customer.ID,
		Order_date:     time.Now(),
		Total_price:    totalPrice,
		Payment_status: false,
	}

	if err := config.DB.Save(&transaction).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": err.Error(),
		})
	}

	for _, td := range req {
		var product models.Product
		if err := config.DB.Model(&models.Product{}).Where("id = ?", td.Product_id).Where("is_verified = true").First(&product).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
				"error": err.Error(),
			})
		}

		transaction_details := models.Transaction_Detail{
			Product_id:     td.Product_id,
			Transaction_id: transaction.ID,
			Seller_id:      product.Seller_id,
			Qty:            td.Qty,
			Price:          product.Price,
		}

		log.Println(transaction_details)

		if err := config.DB.Create(&transaction_details).Error; err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
				"message": "Failed to create product transaction detail",
			})
		}

	}

	return nil
}

func (p *ProductRepository) CreateProductReview(user_id uint, req models.Product_Review_Request) error {
	var customer models.Customer
	if err := config.DB.Model(&models.Customer{}).Where("user_id = ?", user_id).First(&customer).Error; err != nil {
		return err
	}

	var transaction_detail models.Transaction_Detail
	if err := config.DB.Model(&models.Transaction_Detail{}).Where("product_id = ?", req.Product_id).
		Joins("JOIN transactions ON transactions.id = transaction_details.transaction_id").
		Where("cust_id = ?", customer.ID).Where("payment_status = true").First(&transaction_detail).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Customer has not purchased this product or the payment is not completed yet",
			"error":   err.Error(),
		})
	}

	log.Println(transaction_detail)

	review := models.Product_Review{
		Product_id:  req.Product_id,
		Cust_id:     customer.ID,
		Review:      req.Review,
		Review_date: time.Now(),
	}

	log.Println(review)

	if err := config.DB.Save(&review).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, echo.Map{
			"message": "Failed to add review",
		})
	}

	return nil

}
