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
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)

	user := models.User{
		Name:     data["name"],
		Email:    data["email"],
		Password: password,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "failed to create user",
		})
	}

	token, err := GenerateJWT(user)
	fmt.Println("Signup token", token)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to login",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(user)
}



func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to login",
		})
	}

	var user models.User
	if err := database.DB.Where("email = ?", data["email"]).First(&user).Error; err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "failed to find user with corresponding email",
		})
	}

	if user.Id == 0 {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}

	token, err := GenerateJWT(user)
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "failed to login",
		})
	}

	cookie := fiber.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(5 * time.Minute),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	return c.JSON(user)
}



func User(c *fiber.Ctx) error {
	claims, err := AuthenticateUser(c)
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claimsValue := *claims
	userEmail := claimsValue["user_email"]
	
	var user models.User
	if err = database.DB.Where("email = ?", userEmail).First(&user).Error; err != nil {
		return c.JSON(fiber.Map{
			"message": "failed to create user",
		})
	}

	return c.JSON(user)
}



func Logout(c *fiber.Ctx) error {
	ExpireToken(c)

	return c.JSON(fiber.Map{
		"message": "success",
	})
}