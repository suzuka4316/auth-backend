package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/database"
	"github.com/suzuka4316/auth-backend/routes"
)


func init() {
	if err := godotenv.Load(); err != nil {
		panic(".env file found")
	}
}

func main() {
	database.Connect()

	app := fiber.New()

	// accept request from frontend
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
	}))

	routes.Setup(app)

	app.Listen(":8000")
}