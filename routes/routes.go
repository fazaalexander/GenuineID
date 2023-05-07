package routes

import (
	"github.com/fazaalexander/GenuineID/controllers"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func New() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Recover())

	e.POST("/login", controllers.Login)
	e.POST("/register", controllers.Register)
	return e
}
