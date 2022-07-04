package controllers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/suzuka4316/auth-backend/database"
	"github.com/suzuka4316/auth-backend/models"
	"golang.org/x/crypto/bcrypt"
)


func Signup(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(models.Response {
			Success: false,
			Message: "Signup:: failed to parse request body",
		})
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}
	if err := user.Validate(); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(models.Response {
			Success: false,
			Message: fmt.Sprintf("Signup:: user validation error: %v", err),
		})
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(models.Response {
			Success: false,
			Message: "Signup:: failed to create user",
		})
	}

	token, err := GenerateJWT(user)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(models.Response {
			Success: false,
			Message: "Signup:: failed to generate token",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(user)
}



func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(models.Response {
			Success: false,
			Message: "Login:: failed to parse request body",
		})
	}

	var user models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(models.Response{
			Success: false,
			Message: "Login:: failed to find user with corresponding email",
		})
	}

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(models.Response{
			Success: false,
			Message: "Login:: user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(models.Response{
			Success: false,
			Message: "Login:: incorrect password",
		})
	}

	token, err := GenerateJWT(user)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(models.Response{
			Success: false,
			Message: "Login:: failed to generate token",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(1 * time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(user)
}



func User(c *fiber.Ctx) error {
	claims, err := AuthenticateUser(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(models.Response{
			Success: false,
			Message: "User:: unauthenticated user",
		})
	}

	claimsValue := *claims
	userEmail := claimsValue["user_email"]
	
	var user models.User
	if err = database.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(models.Response{
			Success: false,
			Message: "User:: failed to find user with corresponding email",
		})
	}

	return c.JSON(user)
}



func Logout(c *fiber.Ctx) error {
	ExpireToken(c)

	return c.JSON(models.Response{
		Success: true,
		Message: "Logout:: logout successful",
	})
}