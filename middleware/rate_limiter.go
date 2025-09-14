package middleware

import (
	"net/http"
	"time"

	"github.com/AltSumpreme/Medistream.git/config"
	"github.com/AltSumpreme/Medistream.git/services"
	"github.com/AltSumpreme/Medistream.git/utils"
	"github.com/gin-gonic/gin"
)

func RateLimiterMiddleware() gin.HandlerFunc {
	rl := services.NewRateLimiter(config.Rdb, config.Ctx, 100, time.Minute)
	return func(c *gin.Context) {
		ip := c.ClientIP()

		allowed, err := rl.Allow(ip)

		if err != nil {
			utils.Log.Errorf("RateLimiter: Redis error - %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		if !allowed {
			utils.Log.Warnf("RateLimiter: Too many requests from IP %s", ip)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func StrictRateLimiterMiddleware() gin.HandlerFunc {
	rl := services.NewRateLimiter(config.Rdb, config.Ctx, 10, time.Minute)
	return func(c *gin.Context) {
		ip := c.ClientIP()

		allowed, err := rl.Allow(ip)

		if err != nil {
			utils.Log.Errorf("RateLimiter: Redis error - %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}
		if !allowed {
			utils.Log.Warnf("RateLimiter: Too many requests from IP %s", ip)
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
			c.Abort()
			return
		}
		c.Next()
	}
}
