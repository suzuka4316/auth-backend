package controllers

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/joho/godotenv"
	"github.com/suzuka4316/auth-backend/models"
)

func getAPISecret() string {
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

	return token.SignedString([]byte(getAPISecret()))
}


func AuthenticateUser(r *http.Request) (*jwt.MapClaims, error) {
	cookie, err := r.Cookie("jwt")
	if err != nil {
		return &jwt.MapClaims{}, err
	}

	token, err := jwt.ParseWithClaims(cookie.Value, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(getAPISecret()), err
	})
	if err != nil {
		return &jwt.MapClaims{}, err
	}

	return token.Claims.(*jwt.MapClaims), nil
}
