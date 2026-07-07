package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/config"
	"github.com/testgithubanjali/url-shortener/internal/database"
)

func main() {

	// Load environment variables
	config.LoadConfig()

	// Connect to PostgreSQL
	database.ConnectDB()

	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "URL Shortener API is running",
			"port":    config.AppConfig.Port,
		})
	})

	router.Run(":" + config.AppConfig.Port)
}
