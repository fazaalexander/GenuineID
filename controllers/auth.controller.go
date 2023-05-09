package controllers

import (
	"net/http"

	"github.com/fazaalexander/GenuineID/models"
	"github.com/fazaalexander/GenuineID/services"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtCustomClaims struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func Login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	t, err := services.GetAuthRepository().Login(username, password)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token": t,
	})
}

func Register(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "failed",
		})
	}

	user, err := services.GetAuthRepository().Register(user)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success create new user",
		"user":    user,
	})
}
