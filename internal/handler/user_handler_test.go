package handler

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/raihankhan/jwt-auth-system/config"
)

func TestGenerateJWTToken(t *testing.T) {
	// 1. Set up test environment
	cfg := &config.AppConfig{
		JWT: config.JWT{
			SecretKey:          "test-secret-key", // Use a test secret key
			TokenExpiryMinutes: 15,
		},
	}

	gin.SetMode(gin.TestMode)                             // Set Gin to test mode
	c, _ := gin.CreateTestContext(httptest.NewRecorder()) // Create a test Gin context
	c.Set("config", cfg)                                  // Set the test config in the context

	handler := NewUserHandler(nil) // We don't need DB for this test

	userID := uint(123) // Example user ID

	// 2. Call the function to be tested
	tokenString, err := handler.generateJWTToken(userID, c)

	// 3. Assertions - Check for expected outcomes
	if err != nil {
		t.Errorf("generateJWTToken() error = %v, wantErr nil", err) // Check for no error
		return
	}
	if tokenString == "" {
		t.Errorf("generateJWTToken() tokenString is empty, want not empty") // Check for token string not empty
		return
	}

	// (You could add more detailed JWT verification here if needed, like decoding and checking claims)
	// For this basic example, we are just checking for no error and a non-empty token string.

	t.Logf("JWT Token generated successfully: %s", tokenString) // Log success message
}
