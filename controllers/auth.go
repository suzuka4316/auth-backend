package controllers

import (
	"fmt"
	"os"

	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/models"
)

func getAPISecret () string {
	err := godotenv.Load()
	if err != nil {
		panic("getAPISecret:: Error loading .env file")
	} else {
		fmt.Println("getAPISecret:: getting the env value")
	}
	return os.Getenv("API_SECRET")
}


func GenerateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["user_email"] = user.Email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix() //Token expires after 1 hour
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(getAPISecret()))

	if err != nil {
		return "", err
	}
	return tokenString, nil
}



func AuthenticateUser(c *fiber.Ctx) (*jwt.MapClaims, error) {
	cookies := c.Cookies("jwt")

	token, err := jwt.ParseWithClaims(cookies, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(getAPISecret()), nil
	})

	if err != nil {
		return &jwt.MapClaims{}, err
	}

	return token.Claims.(*jwt.MapClaims), nil
}



func ExpireToken(c *fiber.Ctx) {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(time.Hour * -1),
		HTTPOnly: true,
	}

	c.Cookie(&cookie)
}