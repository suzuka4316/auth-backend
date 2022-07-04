package tests

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/models"
	"github.com/suzuka4316/auth-backend/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// var userInstance = models.User{}
var TESTDB *gorm.DB

func TestMain(m *testing.M) {
	err := godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Setup() *fiber.App {
	app := fiber.New()

	// // accept request from frontend
	// app.Use(cors.New(cors.Config{
	// 	AllowCredentials: true,
	// }))

	routes.Setup(app)

	return app
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TEST_DB_DRIVER")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_PASSWORD"), os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_NAME"))
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Printf("Cannot connect to %v database\n", DB)
		log.Fatal("This is the error:", err)
	} else {
		TESTDB = DB
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}
}

func refreshUserTable() error {
	err := TESTDB.Migrator().DropTable(&models.User{})
	if err != nil {
		return err
	}
	err = TESTDB.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Name: "setup test",
		Email: "setup@test.com",
		Password: []byte("password"),
	}

	err = TESTDB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

// func seedUsers() ([]models.User, error) {

// 	var err error
// 	if err != nil {
// 		return nil, err
// 	}
// 	users := []models.User{
// 		models.User{
// 			Nickname: "Steven victor",
// 			Email:    "steven@gmail.com",
// 			Password: "password",
// 		},
// 		models.User{
// 			Nickname: "Kenny Morris",
// 			Email:    "kenny@gmail.com",
// 			Password: "password",
// 		},
// 	}

// 	for i, _ := range users {
// 		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
// 		if err != nil {
// 			return []models.User{}, err
// 		}
// 	}
// 	return users, nil
// }

// func refreshUserAndPostTable() error {

// 	err := server.DB.DropTableIfExists(&models.User{}, &models.Post{}).Error
// 	if err != nil {
// 		return err
// 	}
// 	err = server.DB.AutoMigrate(&models.User{}, &models.Post{}).Error
// 	if err != nil {
// 		return err
// 	}
// 	log.Printf("Successfully refreshed tables")
// 	return nil
// }

// func seedOneUserAndOnePost() (models.Post, error) {

// 	err := refreshUserAndPostTable()
// 	if err != nil {
// 		return models.Post{}, err
// 	}
// 	user := models.User{
// 		Nickname: "Sam Phil",
// 		Email:    "sam@gmail.com",
// 		Password: "password",
// 	}
// 	err = server.DB.Model(&models.User{}).Create(&user).Error
// 	if err != nil {
// 		return models.Post{}, err
// 	}
// 	post := models.Post{
// 		Title:    "This is the title sam",
// 		Content:  "This is the content sam",
// 		AuthorID: user.ID,
// 	}
// 	err = server.DB.Model(&models.Post{}).Create(&post).Error
// 	if err != nil {
// 		return models.Post{}, err
// 	}
// 	return post, nil
// }

// func seedUsersAndPosts() ([]models.User, []models.Post, error) {

// 	var err error

// 	if err != nil {
// 		return []models.User{}, []models.Post{}, err
// 	}
// 	var users = []models.User{
// 		models.User{
// 			Nickname: "Steven victor",
// 			Email:    "steven@gmail.com",
// 			Password: "password",
// 		},
// 		models.User{
// 			Nickname: "Magu Frank",
// 			Email:    "magu@gmail.com",
// 			Password: "password",
// 		},
// 	}
// 	var posts = []models.Post{
// 		models.Post{
// 			Title:   "Title 1",
// 			Content: "Hello world 1",
// 		},
// 		models.Post{
// 			Title:   "Title 2",
// 			Content: "Hello world 2",
// 		},
// 	}

// 	for i, _ := range users {
// 		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
// 		if err != nil {
// 			log.Fatalf("cannot seed users table: %v", err)
// 		}
// 		posts[i].AuthorID = users[i].ID

// 		err = server.DB.Model(&models.Post{}).Create(&posts[i]).Error
// 		if err != nil {
// 			log.Fatalf("cannot seed posts table: %v", err)
// 		}
// 	}
// 	return users, posts, nil
// }