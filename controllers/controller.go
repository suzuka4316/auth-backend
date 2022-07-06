package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"net/http"
	"github.com/suzuka4316/auth-backend/models"
	"github.com/suzuka4316/auth-backend/responses"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}



func (s *Server) Signup(w http.ResponseWriter, r *http.Request) {
	var err error
	
	var user *models.User = &models.User{}
	err = json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Validate("")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	userCreated, err := user.SaveUser(s.DB, hashedPassword)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}

	token, err := GenerateJWT(*userCreated)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(1 * time.Hour),
	})

	w.Header().Set("Location", fmt.Sprintf("%s%s/%d", r.Host, r.RequestURI, userCreated.Id))
	responses.JSON(w, http.StatusCreated, userCreated)
}



func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = user.Validate("login")
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	existingUser, err := models.GetUserByEmail(s.DB, user.Email)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	err = VerifyPassword(existingUser.Password, user.Password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}
	
	token, err := GenerateJWT(*existingUser)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name: "jwt",
		Value: token,
		Expires: time.Now().Add(1 * time.Hour),
	})

	responses.JSON(w, http.StatusOK, existingUser)
}



func (s *Server) User(w http.ResponseWriter, r *http.Request) {
	claims, err := AuthenticateUser(r)
	if err != nil {
		responses.JSON(w, http.StatusUnauthorized, err)
		return
	}

	claimsValue := *claims
	userEmail := claimsValue["user_email"]

	existingUser, err := models.GetUserByEmail(s.DB, userEmail)
	if err != nil {
		responses.ERROR(w, http.StatusUnprocessableEntity, err)
		return
	}

	responses.JSON(w, http.StatusOK, existingUser)
}


func (s *Server) Logout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name: "jwt",
		Value: "",
		Expires:  time.Now().Add(time.Hour * -1),
		HttpOnly: true,
	})

	responses.JSON(w, http.StatusOK, "")
}