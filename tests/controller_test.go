package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/suzuka4316/auth-backend/utils"
)

func TestSignUp(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Printf("TestSignUp:: refreshUserTable err: %v", err)
	}

	testcases := []struct {
		bodyJson     string
		statusCode   int
		name         string
		email        string
		errorMessage string
	}{
		{
			bodyJson:     `{"name":"test", "email": "test@test.com", "password": "password"}`,
			statusCode:   201,
			name:         "test",
			email:        "test@test.com",
			errorMessage: "",
		},
		{
			bodyJson:     `{"name":"test2", "email": "test@test.com", "password": "password"}`,
			statusCode:   500,
			errorMessage: utils.EmailTaken,
		},
		{
			bodyJson:     `{"name":"test2", "email": "testtest.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: utils.InvalidEmail,
		},
		{
			bodyJson:     `{"name":"", "email": "test2@test.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: utils.NameRequired,
		},
		{
			bodyJson:     `{"name":"test2", "email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: utils.EmailRequired,
		},
		{
			bodyJson:     `{"name":"test2", "email": "test2@test.com", "password": ""}`,
			statusCode:   422,
			errorMessage: utils.PasswordRequired,
		},
	}

	for _, tc := range testcases {
		req, err := http.NewRequest("POST", "/signup", bytes.NewBufferString(tc.bodyJson))
		if err != nil {
			t.Errorf("TestSignUp():: NewRequest err: %v", err)
		}
		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Signup)
		handler.ServeHTTP(rr, req)

		responseMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
		if err != nil {
			log.Printf("TestSignUp():: Cannot convert to json: %v", err)
		}
		assert.Equal(t, rr.Code, tc.statusCode)

		if tc.statusCode == 201 {
			assert.Equal(t, tc.name, responseMap["name"])
			assert.Equal(t, tc.email, responseMap["email"])
		}

		if tc.statusCode == 422 || tc.statusCode == 500 && tc.errorMessage != "" {
			assert.Equal(t, tc.errorMessage, responseMap["error"])
		}
	}
}


func TestLogin(t *testing.T) {
	err := refreshUserTable()
	if err != nil {
		log.Printf("TestLogin():: refreshUserTable err: %v", err)
	}

	_, err = seedOneUser()
	if err != nil {
		log.Printf("TestLogin():: seedOneUser err: %v", err)
	}

	testcases := []struct {
		bodyJson     string
		statusCode   int
		email        string
		password     string
		errorMessage string
	}{
		{
			bodyJson:     `{"email": "test@test.com", "password": "password"}`,
			statusCode:   200,
			errorMessage: "",
		},
		{
			bodyJson:     `{"email": "test@test.com", "password": "wrong password"}`,
			statusCode:   422,
			errorMessage: utils.IncorrectPassword,
		},
		{
			bodyJson:     `{"email": "test1@test.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: utils.EmailNotFound,
		},
		{
			bodyJson:     `{"email": "testest.com", "password": "password"}`,
			statusCode:   422,
			errorMessage: utils.InvalidEmail,
		},
		{
			bodyJson:     `{"email": "", "password": "password"}`,
			statusCode:   422,
			errorMessage: utils.EmailRequired,
		},
		{
			bodyJson:     `{"email": "test@test.com", "password": ""}`,
			statusCode:   422,
			errorMessage: utils.PasswordRequired,
		},
	}

	for _, tc := range testcases {
		req, err := http.NewRequest("POST", "/login", bytes.NewBufferString(tc.bodyJson))
		if err != nil {
			t.Errorf("TestLogin():: NewRequest err: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(server.Login)
		handler.ServeHTTP(rr, req)
		assert.Equal(t, rr.Code, tc.statusCode)
		if tc.statusCode == 200 {
			assert.NotEqual(t, rr.Body.String(), "")
		}


		if tc.statusCode == 422 && tc.errorMessage != "" {
			responseMap := make(map[string]interface{})
			err = json.Unmarshal([]byte(rr.Body.String()), &responseMap)
			if err != nil {
				t.Errorf("TestLogin():: Cannot convert to json: %v", err)
			}
			
			fmt.Printf("\nrecord code: %v responseMap: %v tc statuscode: %v tc email: %v tc password: %v\n", rr.Code, responseMap["error"], tc.statusCode, tc.email, tc.password)
			assert.Equal(t, tc.errorMessage, responseMap["error"])
		}
	}
}