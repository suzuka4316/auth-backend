package database

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	err := godotenv.Load()
	if err != nil {
		panic("Connect:: Error loading .env file")
	} else {
		fmt.Println("Connect:: getting the env values")
	}

	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	// port := os.Getenv("DB_PORT")
	// host := os.Getenv("DB_HOST")
	container := os.Getenv("DB_CONTAINER_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", user, password, container, dbname)
	connection, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Connect:: could not connect to the database")
	}

	DB = connection

	connection.AutoMigrate(&models.User{})
}