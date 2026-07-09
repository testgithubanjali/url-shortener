package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/testgithubanjali/url-shortener/internal/database"
)

const (
	MaxRequests = 100
	Window      = time.Minute
)

func RateLimit() gin.HandlerFunc {

	return func(c *gin.Context) {

		ip := c.ClientIP()

		key := "rate:" + ip

		count, err := database.RedisClient.Incr(
			database.Ctx,
			key,
		).Result()

		if err != nil {

			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Rate limiter failed",
			})

			c.Abort()
			return
		}

		// First request
		if count == 1 {

			database.RedisClient.Expire(
				database.Ctx,
				key,
				Window,
			)
		}

		if count > MaxRequests {

			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many requests",
			})

			c.Abort()
			return
		}

		c.Next()
	}
}
