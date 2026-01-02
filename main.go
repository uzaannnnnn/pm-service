package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(corsMiddleware())
	r.POST("/api/pm", handlePM)
	r.GET("/api/news/:age", handleGetNews)
	r.GET("/api/news/:age/:id", handleGetNewsDetail)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
