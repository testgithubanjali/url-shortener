package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"github.com/testgithubanjali/url-shortener/internal/config"
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

	// Get logged-in user ID from JWT
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
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

		err = database.DB.
			Where("short_code = ?", code).
			First(&existingURL).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
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

	// Calculate expiration date
	var expiresAt *time.Time

	if req.ExpiresInDays > 0 {
		t := time.Now().AddDate(0, 0, req.ExpiresInDays)
		expiresAt = &t
	}

	url := models.URL{
		UserID:      userUUID,
		OriginalURL: req.OriginalURL,
		ShortCode:   shortCode,
		ExpiresAt:   expiresAt,
	}

	if err := database.DB.Create(&url).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to save URL",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "URL shortened successfully",
		"short_url": "http://localhost:" + config.AppConfig.Port + "/" + shortCode,
	})
}

func RedirectURL(c *gin.Context) {

	shortCode := c.Param("shortCode")

	cacheKey := "url:" + shortCode

	// Check Redis first
	originalURL, err := database.RedisClient.Get(database.Ctx, cacheKey).Result()

	if err == nil {

		log.Println("✅ Cache HIT:", shortCode)

		database.DB.Model(&models.URL{}).
			Where("short_code = ?", shortCode).
			Update("click_count", gorm.Expr("click_count + 1"))

		c.Redirect(http.StatusFound, originalURL)
		return
	}

	if err != nil && !errors.Is(err, redis.Nil) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Redis error",
		})
		return
	}

	log.Println("❌ Cache MISS:", shortCode)

	var url models.URL

	err = database.DB.
		Where("short_code = ?", shortCode).
		First(&url).Error

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Short URL not found",
		})
		return
	}

	// Check if URL has expired
	if url.ExpiresAt != nil && time.Now().After(*url.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{
			"error": "This short URL has expired",
		})
		return
	}

	// Save to Redis for 10 minutes
	err = database.RedisClient.Set(
		database.Ctx,
		cacheKey,
		url.OriginalURL,
		10*time.Minute,
	).Err()

	if err != nil {
		log.Println("Failed to cache URL:", err)
	}

	// Increase click count
	database.DB.Model(&url).
		Update("click_count", gorm.Expr("click_count + 1"))

	c.Redirect(http.StatusFound, url.OriginalURL)
}
