package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/testgithubanjali/url-shortener/internal/database"
	"github.com/testgithubanjali/url-shortener/internal/dto"
	"github.com/testgithubanjali/url-shortener/internal/models"
	"github.com/testgithubanjali/url-shortener/internal/services"
)

func ShortenURL(c *gin.Context) {

	var req dto.ShortenURLRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get logged-in user's ID from JWT middleware
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User not found",
		})
		return
	}

	// Generate a unique short code
	var shortCode string

	for {
		code, err := services.GenerateShortCode()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate short code",
			})
			return
		}

		var existingURL models.URL

		err = database.DB.Where("short_code = ?", code).First(&existingURL).Error

		if err != nil {
			shortCode = code
			break
		}
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	url := models.URL{
		UserID:      userUUID,
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
	}

	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save URL",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "URL shortened successfully",
		"short_url": "http://localhost:8080/" + shortCode,
	})
}
func RedirectURL(c *gin.Context) {

	shortCode := c.Param("shortCode")

	var url models.URL

	err := database.DB.Where("short_code = ?", shortCode).First(&url).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	// Increase click count
	url.ClickCount++

	database.DB.Save(&url)

	// Redirect user
	c.Redirect(http.StatusFound, url.OriginalURL)
}
