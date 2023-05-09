package config

import (
	"fmt"
	"log"
	"os"

	"github.com/fazaalexander/GenuineID/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB *gorm.DB
)

func init() {
	InitDB()
	InitialMigration()
}

type Config struct {
	DB_Username string
	DB_Password string
	DB_Port     string
	DB_Host     string
	DB_Name     string
}

func InitDB() {
	erro := godotenv.Load()
	if erro != nil {
		log.Fatal("Error loading .env file")
	}

	db_username := os.Getenv("DB_Username")
	db_password := os.Getenv("DB_Password")
	db_port := os.Getenv("DB_Port")
	db_host := os.Getenv("DB_Host")
	db_name := os.Getenv("DB_Name")

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		db_username,
		db_password,
		db_host,
		db_port,
		db_name,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		panic(err)
	}
}

func InitialMigration() {
	DB.AutoMigrate(
		&models.Admin{},
		&models.Seller{},
		&models.Customer{},
		&models.Product{},
		&models.Product_Auth{},
		&models.Product_Type{},
		&models.Product_Review{},
		&models.Transaction{},
		&models.Transaction_Detail{},
		// &models.Payment{},
	)
}
