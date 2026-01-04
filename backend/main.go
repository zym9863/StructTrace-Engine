package main

import (
	"gin/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// API v1 routes
	api := r.Group("/api/v1")
	{
		// Data structure operations
		api.POST("/operations", handlers.HandleOperation)
		api.POST("/reset", handlers.HandleReset)

		// Benchmark endpoints
		api.POST("/benchmark/start", handlers.HandleBenchmarkSSE)
		api.POST("/benchmark/stop", handlers.HandleStopBenchmark)
		api.GET("/benchmark/status", handlers.HandleBenchmarkStatus)

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "ok",
				"service": "StructTrace Engine API",
			})
		})
	}

	r.Run(":8080")
}
