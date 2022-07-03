package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/suzuka4316/auth-backend/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/signup", controllers.Signup)
	app.Post("/login", controllers.Login)
	app.Get("/user", controllers.User)
	app.Post("/logout", controllers.Logout)
}