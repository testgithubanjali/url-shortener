package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/handlers"
)

func SetupRoutes(router *gin.Engine) {

	router.POST("/register", handlers.Register)
}
