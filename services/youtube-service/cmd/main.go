package main

import (
	"log"
	"os"
	"youtube-service/internal/handler"
	"youtube-service/internal/repository"
	"youtube-service/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8081"
	}

	// Initialize database
	db, err := repository.NewPostgresDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	videoRepo := repository.NewVideoRepository(db)

	// Initialize service
	youtubeService := service.NewYouTubeService(videoRepo)

	// Initialize handler
	youtubeHandler := handler.NewYouTubeHandler(youtubeService)

	router := gin.Default()

	// Routes
	router.GET("/health", youtubeHandler.HealthCheck)
	router.POST("/api/extract", youtubeHandler.ExtractVideo)
	router.GET("/api/video/:id", youtubeHandler.GetVideo)
	router.GET("/api/subtitles/:id", youtubeHandler.GetSubtitles)

	log.Printf("YouTube service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start YouTube service: %v", err)
	}
}