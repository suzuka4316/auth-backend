package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/controllers"
	"github.com/suzuka4316/auth-backend/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var server = controllers.Server{}

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("../.env"))
	if err != nil {
		log.Fatalf("dbconnection_test TestMain:: Error getting .env values %v\n", err)
	}

	Database()

	os.Exit(m.Run())
}


func Database() {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))

	server.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("dbconnection_test Database():: Failed to connect to database err: %v", err)
	} else {
		log.Printf("dbconnection_test Database():: connected to the database\n")
	}
}


func refreshUserTable() error {
	err := server.DB.Migrator().DropTable(&models.User{})
	if err != nil {
		return err
	}
	
	err = server.DB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	log.Printf("dbconnection_test refreshUserTable():: refreshed user table")
	return nil
}


func seedOneUser() (models.User, error) {
	err := refreshUserTable()
	if err != nil {
		log.Printf("dbconnection_test seedOneUser():: refreshUserTable err: %v", err)
		return models.User{}, err
	}

	// need to hash password before because plain password is never saved into database
	hashedPassword, err := controllers.HashPassword("password")
	if err != nil {
		log.Printf("dbconnection_test seedOneUser():: HashPassword err: %v", err)
		return models.User{}, err
	}

	user := models.User{
		Name:     "test",
		Email:    "test@test.com",
		Password: string(hashedPassword),
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		log.Printf("dbconnection_test seedOneUser():: saving user err: %v", err)
		return models.User{}, err
	}

	return user, nil
}