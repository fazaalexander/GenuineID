package services

import (
	"net/http"

	"github.com/fazaalexander/GenuineID/config"
	"github.com/fazaalexander/GenuineID/models"
	"github.com/labstack/echo/v4"
)

// Interface Product
type IUserService interface {
	ResetPassword(user_id uint, hashedPassword string) error
}

type UserRepository struct {
	Func IUserService
}

var userRepository IUserService

func init() {
	ur := &UserRepository{}
	ur.Func = ur

	userRepository = ur
}

func GetUserRepository() IUserService {
	return userRepository
}

func SetUserRepository(ur IUserService) {
	userRepository = ur
}

func (u *UserRepository) ResetPassword(user_id uint, hashedPassword string) error {
	if err := config.DB.Model(&models.User{}).Where("id = ?", user_id).Update("password", hashedPassword).Error; err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, echo.Map{
			"message": "Reset Password Failed",
		})
	}

	return nil
}
