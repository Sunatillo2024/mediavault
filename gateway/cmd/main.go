package main

import (
	"log"
	"os"

	"gateway/internal/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}

	youtubeServiceURL := os.Getenv("YOUTUBE_SERVICE_URL")
	if youtubeServiceURL == "" {
		youtubeServiceURL = "http://localhost:8081"
	}

	instagramServiceURL := os.Getenv("INSTAGRAM_SERVICE_URL")
	if instagramServiceURL == "" {
		instagramServiceURL = "http://localhost:8082"
	}

	router := gin.Default()

	// CORS configuration
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Initialize handler
	gatewayHandler := handler.NewGatewayHandler(youtubeServiceURL, instagramServiceURL)

	// Routes
	router.GET("/health", gatewayHandler.HealthCheck)
	router.POST("/api/extract", gatewayHandler.ExtractContent)
	router.GET("/api/content/:id", gatewayHandler.GetContent)
	router.POST("/api/youtube/extract", gatewayHandler.ExtractYouTube)
	router.GET("/api/youtube/video/:id", gatewayHandler.GetYouTubeVideo)
	router.POST("/api/instagram/extract", gatewayHandler.ExtractInstagram)
	router.GET("/api/instagram/post/:id", gatewayHandler.GetInstagramPost)

	log.Printf("Gateway service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start gateway: %v", err)
	}
}