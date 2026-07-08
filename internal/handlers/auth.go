package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/database"
	"github.com/testgithubanjali/url-shortener/internal/dto"
	"github.com/testgithubanjali/url-shortener/internal/models"
	"github.com/testgithubanjali/url-shortener/internal/utils"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest

	log.Printf("Registration attempt for email %s", req.Email)

	// Read JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Registration validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if the email already exists
	var existingUser models.User

	err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		log.Printf("Registration failed: email %s already exists", req.Email)
		c.JSON(http.StatusConflict, gin.H{
			"error": "Email already exists",
		})
		return
	}

	// Hash the password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		log.Printf("Password hashing failed for email %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}

	// Create user object
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save user to database
	err = database.DB.Create(&user).Error
	if err != nil {
		log.Printf("User creation failed for email %s: %v", req.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user",
		})
		return
	}

	log.Printf("User registered successfully: %s", user.Email)
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}

func Login(c *gin.Context) {
	var req dto.LoginRequest

	// Read request body
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("Login validation failed: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Login attempt for email %s", req.Email)

	// Find user by email
	var user models.User

	err := database.DB.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		log.Printf("Login failed: user not found for email %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Compare password
	err = utils.CheckPassword(user.Password, req.Password)
	if err != nil {
		log.Printf("Login failed: invalid password for email %s", req.Email)
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid email or password",
		})
		return
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID.String())
	if err != nil {
		log.Printf("Token generation failed for user %s: %v", user.Email, err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	log.Printf("Login successful for user %s", user.Email)
	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
	})
}

func Profile(c *gin.Context) {
	userID, _ := c.Get("user_id")
	log.Printf("Profile requested for user %v", userID)

	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome!",
		"user_id": userID,
	})
}
