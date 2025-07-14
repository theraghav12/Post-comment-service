package main

import (
	"log"
	"net/http"
	"time"

	"post-comments-api/config"
	"post-comments-api/routes"
	"post-comments-api/utils"

	"github.com/gin-gonic/gin"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	// Load configuration
	config.LoadConfig()

	// Initialize database connection
	utils.InitDB()

	// Initialize routes
	r := routes.SetupRouter()

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Start server
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on :%s", port)
	r.Run(":" + port)
}
