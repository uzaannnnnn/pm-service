package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func corsMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{"http://localhost:3000"}
	if envOrigins := strings.TrimSpace(os.Getenv("CORS_ORIGINS")); envOrigins != "" {
		allowedOrigins = strings.Split(envOrigins, ",")
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			for _, allowed := range allowedOrigins {
				if strings.TrimSpace(allowed) == origin {
					c.Header("Access-Control-Allow-Origin", origin)
					break
				}
			}
		}
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
