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
	e.PUT("/resetpassword", controllers.ResetPassword)

	productGroup := e.Group("/products", middlewares.SellerTokenVerify)
	productGroup.POST("/", controllers.CreateProduct)
	productGroup.GET("/", controllers.GetAllProducts)
	productGroup.GET("/:id", controllers.GetProductByID)
	productGroup.DELETE("/:id", controllers.DeleteProduct)
	productGroup.PUT("/:id", controllers.UpdateProduct)

	productAuthGroup := e.Group("/products", middlewares.AdminTokenVerify)
	productAuthGroup.PUT("/authenticate", controllers.AuthenticateProduct)

	custProductGroup := e.Group("/products", middlewares.CustomerTokenVerify)
	custProductGroup.GET("/search", controllers.SearchProductByName)
	custProductGroup.GET("/search/:id", controllers.SearchProductByID)
	custProductGroup.GET("/category/search", controllers.SearchProductByType)
	custProductGroup.POST("/checkout", controllers.ProductCheckout)
	custProductGroup.POST("/review", controllers.CreateProductReview)
	return e
}
