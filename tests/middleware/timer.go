package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func RequestTimer() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)
		status := c.Writer.Status()

		log.Printf("Timer %s %s | %d | %v",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
		)
	}
}
