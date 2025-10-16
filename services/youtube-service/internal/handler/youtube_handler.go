package handler

import (
	"net/http"
	"youtube-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type YouTubeHandler struct {
	service *service.YouTubeService
}

type ExtractRequest struct {
	URL               string `json:"url"`
	DownloadVideo     bool   `json:"download_video"`
	DownloadSubtitles bool   `json:"download_subtitles"`
}

func NewYouTubeHandler(service *service.YouTubeService) *YouTubeHandler {
	return &YouTubeHandler{service: service}
}

func (h *YouTubeHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "youtube"})
}

func (h *YouTubeHandler) ExtractVideo(c *gin.Context) {
	var req ExtractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	video, err := h.service.ExtractVideo(req.URL, req.DownloadVideo, req.DownloadSubtitles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (h *YouTubeHandler) GetVideo(c *gin.Context) {
	id := c.Param("id")

	// Validate UUID
	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	video, err := h.service.GetVideo(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	c.JSON(http.StatusOK, video)
}

func (h *YouTubeHandler) GetSubtitles(c *gin.Context) {
	id := c.Param("id")

	if _, err := uuid.Parse(id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid video ID"})
		return
	}

	subtitles, err := h.service.GetSubtitles(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subtitles not found"})
		return
	}

	c.JSON(http.StatusOK, subtitles)
}