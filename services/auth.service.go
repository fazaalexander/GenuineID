package services

import (
	"log"
	"net/http"

	"github.com/fazaalexander/GenuineID/config"
	"github.com/fazaalexander/GenuineID/middlewares"
	"github.com/fazaalexander/GenuineID/models"
	"github.com/fazaalexander/GenuineID/utils"

	"github.com/labstack/echo/v4"
)

// Interface Product
type IAuthService interface {
	Login(username string, password string) (interface{}, error)
	Register(user models.User) (models.User, error)
}

type AuthRepository struct {
	Func IAuthService
}

var authRepository IAuthService

func init() {
	ar := &AuthRepository{}
	ar.Func = ar

	authRepository = ar
}

func GetAuthRepository() IAuthService {
	return authRepository
}

func SetAuthRepository(ar IAuthService) {
	authRepository = ar
}

func (a *AuthRepository) Login(username string, password string) (interface{}, error) {
	var user models.User
	if err := config.DB.Model(&user).Where("username = ?", username).First(&user).Error; err != nil {
		log.Println("username tidak ditemukan")
		return nil, echo.ErrUnauthorized
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		log.Println("password salah")
		return nil, echo.ErrUnauthorized
	}

	return middlewares.CreateToken(&user, user.ID, user.Username, user.Email, user.Role)
}

func (a *AuthRepository) Register(user models.User) (models.User, error) {
	if err := config.DB.Save(&user).Error; err != nil {
		return user, echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	switch user.Role {
	case "customer":
		customer := models.Customer{
			User_id: user.ID,
		}

		if err := config.DB.Save(&customer).Error; err != nil {
			return user, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	case "seller":
		seller := models.Seller{
			User_id: user.ID,
		}

		if err := config.DB.Save(&seller).Error; err != nil {
			return user, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	case "admin":
		admin := models.Admin{
			User_id: user.ID,
		}

		if err := config.DB.Save(&admin).Error; err != nil {
			return user, echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	}

	return user, nil
}
