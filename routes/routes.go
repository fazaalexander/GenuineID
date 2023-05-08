package routes

import (
	"github.com/fazaalexander/GenuineID/controllers"
	"github.com/fazaalexander/GenuineID/middlewares"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())

	e.POST("/login", controllers.Login)
	e.POST("/register", controllers.Register)
	e.PUT("/profile", controllers.Profile)

	productGroup := e.Group("/products", middlewares.SellerTokenVerify)
	productGroup.POST("/", controllers.CreateProduct)
	productGroup.GET("/", controllers.GetAllProducts)
	productGroup.GET("/:id", controllers.GetProductByID)
	productGroup.DELETE("/:id", controllers.DeleteProduct)
	productGroup.PUT("/:id", controllers.UpdateProduct)

	return e
}
