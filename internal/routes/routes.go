package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/handlers"
	"github.com/testgithubanjali/url-shortener/internal/middleware"
)

func SetupRoutes(router *gin.Engine) {

	// Public Routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)

	// Protected Routes
	authorized := router.Group("/")

	authorized.Use(
		middleware.RateLimit(),
		middleware.JWTAuth(),
	)

	{
		authorized.GET("/profile", handlers.Profile)
		authorized.POST("/shorten", handlers.ShortenURL)
		authorized.GET("/analytics/:shortCode", handlers.GetAnalytics)
	}

	// Public Redirect Route (also rate limited)
	router.GET(
		"/:shortCode",
		middleware.RateLimit(),
		handlers.RedirectURL,
	)
}
