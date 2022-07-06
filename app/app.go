package app

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/controllers"
)

var server = controllers.Server{}

func init() {
	if err := godotenv.Load(); err != nil {
		panic("init:: .env file found")
	}
}

func Run() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Run:: Error getting env, %v", err)
	} else {
		fmt.Println("Run:: getting the env values")
	}

	server.Connect(os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_PORT"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))
	
	server.Run(":8000")
}