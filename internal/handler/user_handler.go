package handler

import (
	"github.com/raihankhan/jwt-auth-system/config"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/raihankhan/jwt-auth-system/internal/user" // Adjust import path if needed
)

// LoginRequest defines the request body for user login.
type LoginRequest struct {
	UsernameOrEmail string `json:"usernameOrEmail" binding:"required"` // Allow login with username or email
	Password        string `json:"password" binding:"required"`
}

// LoginUser handles the user login request and generates a JWT.
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Find user by username or email
	var foundUser user.User
	result := h.db.Where("username = ? OR email = ?", req.UsernameOrEmail, req.UsernameOrEmail).First(&foundUser)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // User not found
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user"}) // Database error
		}
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"}) // Password mismatch
		return
	}

	// Generate JWT token
	token, err := h.generateJWTToken(foundUser.ID, c) // Call JWT generation function
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token}) // Return JWT on success
}

// generateJWTToken is a helper function to generate a JWT token for a user.
func (h *UserHandler) generateJWTToken(userID uint, c *gin.Context) (string, error) {
	config := c.MustGet("config").(*config.AppConfig) // Retrieve config from context

	// Create JWT token claims
	claims := jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(time.Minute * time.Duration(config.JWT.TokenExpiryMinutes)).Unix(), // Token expiration time
		"issuedAt":  time.Now().Unix(),
		"issuer":    "jwt-auth-system", // You can customize issuer
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret key
	secretKey := []byte(config.JWT.SecretKey) // Get secret key from config
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// RegisterUserRequest defines the request body for user registration.
type RegisterUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"omitempty,email"` // Optional, but must be email format
	FullName string `json:"fullName"`
}

// UserHandler handles user-related API requests.
type UserHandler struct {
	db *gorm.DB
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

// RegisterUser handles the user registration request.
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req RegisterUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Create User object
	newUser := user.User{
		Username:     req.Username,
		PasswordHash: string(hashedPassword), // Store the hashed password
		Email:        req.Email,
		FullName:     req.FullName,
	}

	// Save user to database
	result := h.db.Create(&newUser)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully", "userID": newUser.ID})
}

// ProtectedEndpoint is a sample protected endpoint that requires authentication.
func (h *UserHandler) ProtectedEndpoint(c *gin.Context) {
	userID := c.MustGet("userID").(uint) // Retrieve userID from context (set by middleware)

	c.JSON(http.StatusOK, gin.H{"message": "Protected endpoint accessed!", "userID": userID})
}
