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
