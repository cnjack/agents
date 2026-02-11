package api

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLogger logs incoming requests
func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(startTime)
		status := c.Writer.Status()

		log.Printf("[%s] %s %d %v",
			c.Request.Method,
			c.Request.URL.Path,
			status,
			latency,
		)
	}
}

// Recovery handles panics and returns 500 error
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic recovered: %v", err)
				c.JSON(500, gin.H{
					"error": "Internal server error",
				})
				c.Abort()
			}
		}()
		c.Next()
	}
}

// RateLimiter limits requests per IP (simple implementation)
func RateLimiter() gin.HandlerFunc {
	// In production, use a proper rate limiter like redis-based
	return func(c *gin.Context) {
		// For now, just pass through
		c.Next()
	}
}

// AuthMiddleware validates authentication (for future use)
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// For now, no authentication required
		// In production, validate API key or JWT token
		c.Next()
	}
}
