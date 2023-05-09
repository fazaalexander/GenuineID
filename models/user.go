package models

import (
	"github.com/fazaalexander/GenuineID/utils"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `json:"username" form:"username" gorm:"uniqueIndex"`
	Email        string `json:"email" form:"email" gorm:"uniqueIndex"`
	Phone_number string `json:"phone_number" form:"phone_number"`
	Address      string `json:"address" form:"address"`
	Password     string `json:"password" form:"password"`
	Role         string `json:"role" form:"role"`
}

func (u *User) BeforeCreate(tc *gorm.DB) (err error) {
	hashPassword, err := utils.HashPassword(u.Password)
	if err != nil {
		return err
	}
	u.Password = hashPassword

	return
}
