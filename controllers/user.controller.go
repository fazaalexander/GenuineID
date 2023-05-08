package controllers

import (
	"net/http"
	"os"
	"strings"

	"github.com/fazaalexander/GenuineID/config"
	"github.com/fazaalexander/GenuineID/models"
	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

func Profile(c echo.Context) error {
	token := c.Request().Header.Get(("Authorization"))
	token = strings.Replace(token, "Bearer ", "", 1)

	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{
			"error": "Invalid token",
		})
	}

	claims := jwtToken.Claims.(*jwtCustomClaims)
	id := claims.ID
	role := claims.Role

	switch role {
	case "customer":
		var request struct {
			Username     string          `json:"username" form:"username" gorm:"uniqueIndex"`
			Email        string          `json:"email" form:"email" gorm:"uniqueIndex"`
			Phone_number string          `json:"phone_number" form:"phone_number"`
			Address      string          `json:"address" form:"address"`
			Customer     models.Customer `json:"detail" form:"detail"`
		}

		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "failed",
			})
		}
		user := models.User{
			Username:     request.Username,
			Email:        request.Email,
			Phone_number: request.Phone_number,
			Address:      request.Address,
		}

		customer := models.Customer{
			FullName: request.Customer.FullName,
		}

		if err := config.DB.Model(&models.User{}).Where("id = ?", id).Updates(models.User{Username: user.Username, Email: user.Email, Phone_number: user.Phone_number, Address: user.Address}).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := config.DB.Model(&models.Customer{}).Where("user_id = ?", id).Updates(models.Customer{FullName: customer.FullName}).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	case "seller":
		var request struct {
			Username     string        `json:"username" form:"username" gorm:"uniqueIndex"`
			Email        string        `json:"email" form:"email" gorm:"uniqueIndex"`
			Phone_number string        `json:"phone_number" form:"phone_number"`
			Address      string        `json:"address" form:"address"`
			Seller       models.Seller `json:"detail" form:"detail"`
		}

		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "failed",
			})
		}
		user := models.User{
			Username:     request.Username,
			Email:        request.Email,
			Phone_number: request.Phone_number,
			Address:      request.Address,
		}

		seller := models.Seller{
			Store_name: request.Seller.Store_name,
		}

		if err := config.DB.Model(&models.User{}).Where("id = ?", id).Updates(models.User{Username: user.Username, Email: user.Email, Phone_number: user.Phone_number, Address: user.Address}).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		if err := config.DB.Model(&models.Seller{}).Where("user_id = ?", id).Updates(models.Seller{Store_name: seller.Store_name}).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	case "admin":
		var request struct {
			Username     string `json:"username" form:"username" gorm:"uniqueIndex"`
			Email        string `json:"email" form:"email" gorm:"uniqueIndex"`
			Phone_number string `json:"phone_number" form:"phone_number"`
			Address      string `json:"address" form:"address"`
		}

		if err := c.Bind(&request); err != nil {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "failed",
			})
		}
		user := models.User{
			Username:     request.Username,
			Email:        request.Email,
			Phone_number: request.Phone_number,
			Address:      request.Address,
		}

		if err := config.DB.Model(&models.User{}).Where("id = ?", id).Updates(models.User{Username: user.Username, Email: user.Email, Phone_number: user.Phone_number, Address: user.Address}).Error; err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	default:
		return echo.NewHTTPError(http.StatusBadRequest, "Unknown user role")
	}

	return c.JSON(http.StatusOK, echo.Map{
		"message": "success update user profile",
	})
}
