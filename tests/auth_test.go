package tests

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suzuka4316/auth-backend/controllers"
)

func TestGenerateJWT(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user, err := seedOneUser()
	if err != nil {
		fmt.Printf("TestSignIn:: seedOneUser error: %v", err)
	}

	token, _ := controllers.GenerateJWT(user)
	assert.NotEqual(t, token, "")
}