package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/config"
	"github.com/testgithubanjali/url-shortener/internal/database"
	"github.com/testgithubanjali/url-shortener/internal/routes"
)

func main() {
	log.Println("Starting URL shortener service")

	// Load environment variables
	config.LoadConfig()
	log.Printf("Configuration loaded for port %s", config.AppConfig.Port)

	// Connect to PostgreSQL
	log.Println("Connecting to PostgreSQL")
	database.ConnectDB()
	log.Println("Connecting to Redis")
	database.ConnectRedis()

	// Create Gin router
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Authorization",
		},
		ExposeHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Register all application routes
	log.Println("Registering application routes")
	routes.SetupRoutes(router)

	// Health Check API
	router.GET("/health", func(c *gin.Context) {
		log.Println("Health check requested")
		c.JSON(http.StatusOK, gin.H{
			"message": "URL Shortener API is running",
			"port":    config.AppConfig.Port,
		})
	})

	// Start the server
	log.Printf("Server listening on port %s", config.AppConfig.Port)
	if err := router.Run(":" + config.AppConfig.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
