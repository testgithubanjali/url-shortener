package main

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/config"
	"github.com/testgithubanjali/url-shortener/internal/database"
	"github.com/testgithubanjali/url-shortener/internal/routes"
)

func main() {

	// Load environment variables
	config.LoadConfig()

	// Connect to PostgreSQL
	database.ConnectDB()
	database.ConnectRedis()

	// Create Gin router
	router := gin.Default()

	// Register all application routes
	routes.SetupRoutes(router)

	// Health Check API
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "URL Shortener API is running",
			"port":    config.AppConfig.Port,
		})
	})

	// Start the server
	router.Run(":" + config.AppConfig.Port)
}
