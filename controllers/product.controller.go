package controllers

import (
	"net/http"
	"strconv"

	"github.com/fazaalexander/GenuineID/middlewares"
	"github.com/fazaalexander/GenuineID/models"
	"github.com/fazaalexander/GenuineID/services"

	"github.com/labstack/echo/v4"
)

// Mendaftarkan barang dagangan
func CreateProduct(c echo.Context) error {
	token := c.Request().Header.Get(("Authorization"))

	claims, err := middlewares.GetClaims(token)
	if err != nil {
		return err
	}

	seller_id := claims.ID

	var req struct {
		Product_type_id uint    `json:"product_type_id" form:"product_type_id"`
		Name            string  `json:"name" form:"name"`
		Description     string  `json:"description" form:"description"`
		Price           float64 `json:"price" form:"price"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	product := models.Product{
		Seller_id:       seller_id,
		Product_type_id: req.Product_type_id,
		Name:            req.Name,
		Description:     req.Description,
		Price:           req.Price,
		Is_Verified:     false,
	}

	if err := services.GetProductRepository().CreateProduct(product); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully create product",
	})
}

// Melihat daftar barang dagangan
func GetAllProducts(c echo.Context) error {
	var products *[]models.ProductResponse

	products, err := services.GetProductRepository().GetProducts()
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":  "success get all products",
		"products": products,
	})
}

// Melihat barang dagangan berdasarkan ID
func GetProductByID(c echo.Context) error {
	var product *models.ProductResponse
	id := c.Param("id")

	product, err := services.GetProductRepository().GetProductByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":  "success get product by id",
		"products": product,
	})
}

// Menghapus barang dagangan
func DeleteProduct(c echo.Context) error {
	id := c.Param("id")

	err := services.GetProductRepository().DeleteProduct(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully delete product",
	})
}

// Mengubah detail barang dagangan
func UpdateProduct(c echo.Context) error {
	id := c.Param("id")
	name := c.FormValue("name")
	description := c.FormValue("description")
	price, ok := strconv.ParseFloat(c.FormValue("price"), 64)
	if ok != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Invalid price",
		})
	}

	if err := services.GetProductRepository().UpdateProduct(id, name, description, price); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully update product data",
	})
}

// Autentikasi barang
func AuthenticateProduct(c echo.Context) error {
	var product_auth *models.Product_Auth

	token := c.Request().Header.Get(("Authorization"))

	claims, err := middlewares.GetClaims(token)
	if err != nil {
		return err
	}

	admin_id := claims.ID

	if err := c.Bind(&product_auth); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"error": err.Error(),
		})
	}

	product_auth, err = services.GetProductRepository().AuthenticateProduct(admin_id, product_auth)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully authenticate product",
		"product": product_auth,
	})
}

// Mencari barang berdasarkan nama
func SearchProductByName(c echo.Context) error {
	var products *[]models.ProductResponse
	name := c.QueryParam("name")

	products, err := services.GetProductRepository().SearchProductByName(name)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":  "Product Found",
		"products": products,
	})
}

// Mencari barang berdasarkan ID
func SearchProductByID(c echo.Context) error {
	var product *models.ProductResponse
	id := c.Param("id")

	product, err := services.GetProductRepository().SearchProductByID(id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Product Found",
		"product": product,
	})
}

// Mencari barang berdasarkan type (kategori)
func SearchProductByType(c echo.Context) error {
	var products *[]models.ProductResponse
	product_type_id := c.QueryParam("type")

	products, err := services.GetProductRepository().SearchProductByType(product_type_id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message":  "Product Found",
		"products": products,
	})
}

// Membeli barang
func ProductCheckout(c echo.Context) error {
	token := c.Request().Header.Get(("Authorization"))

	claims, err := middlewares.GetClaims(token)
	if err != nil {
		return err
	}

	user_id := claims.ID

	var req []models.Transaction_Detail_Request

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	if err := services.GetProductRepository().ProductCheckout(user_id, req); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully create product transaction",
	})
}

// memberikan review
func CreateProductReview(c echo.Context) error {
	token := c.Request().Header.Get(("Authorization"))

	claims, err := middlewares.GetClaims(token)
	if err != nil {
		return err
	}

	user_id := claims.ID

	var req models.Product_Review_Request

	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	err = services.GetProductRepository().CreateProductReview(user_id, req)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "Successfully add review",
	})
}
