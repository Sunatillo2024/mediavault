package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type GatewayHandler struct {
	youtubeServiceURL   string
	instagramServiceURL string
	httpClient          *http.Client
}

type ExtractRequest struct {
	URL                string `json:"url"`
	Platform           string `json:"platform"`
	DownloadVideo      bool   `json:"download_video"`
	DownloadSubtitles  bool   `json:"download_subtitles"`
	DownloadMedia      bool   `json:"download_media"`
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Services  map[string]string `json:"services"`
	Timestamp string            `json:"timestamp"`
}

func NewGatewayHandler(youtubeServiceURL, instagramServiceURL string) *GatewayHandler {
	return &GatewayHandler{
		youtubeServiceURL:   youtubeServiceURL,
		instagramServiceURL: instagramServiceURL,
		httpClient:          &http.Client{},
	}
}

func (h *GatewayHandler) HealthCheck(c *gin.Context) {
	services := make(map[string]string)
	timestamp := "2024-01-01T00:00:00Z" // In real implementation, use actual timestamp

	// Check YouTube service
	youtubeHealth, err := h.checkServiceHealth(h.youtubeServiceURL)
	if err != nil {
		services["youtube"] = "down"
	} else {
		services["youtube"] = youtubeHealth
	}

	// Check Instagram service
	instagramHealth, err := h.checkServiceHealth(h.instagramServiceURL)
	if err != nil {
		services["instagram"] = "down"
	} else {
		services["instagram"] = instagramHealth
	}

	// Determine overall status
	overallStatus := "healthy"
	for _, status := range services {
		if status != "ok" {
			overallStatus = "degraded"
			break
		}
	}

	response := HealthResponse{
		Status:    overallStatus,
		Services:  services,
		Timestamp: timestamp,
	}

	c.JSON(http.StatusOK, response)
}

func (h *GatewayHandler) checkServiceHealth(serviceURL string) (string, error) {
	resp, err := h.httpClient.Get(serviceURL + "/health")
	if err != nil {
		return "down", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "down", fmt.Errorf("service returned status: %d", resp.StatusCode)
	}

	var healthResponse map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&healthResponse); err != nil {
		return "down", err
	}

	return healthResponse["status"], nil
}

func (h *GatewayHandler) ExtractContent(c *gin.Context) {
	var req ExtractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	switch strings.ToLower(req.Platform) {
	case "youtube":
		h.forwardToYouTube(c, "/api/extract", req)
	case "instagram":
		h.forwardToInstagram(c, "/api/extract", req)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported platform"})
	}
}

func (h *GatewayHandler) ExtractYouTube(c *gin.Context) {
	var req ExtractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	h.forwardToYouTube(c, "/api/extract", req)
}

func (h *GatewayHandler) ExtractInstagram(c *gin.Context) {
	var req ExtractRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	h.forwardToInstagram(c, "/api/extract", req)
}

func (h *GatewayHandler) GetContent(c *gin.Context) {
	id := c.Param("id")

	// Try YouTube first
	resp, err := h.httpClient.Get(h.youtubeServiceURL + "/api/video/" + id)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, "application/json", body)
		return
	}

	// Try Instagram
	resp, err = h.httpClient.Get(h.instagramServiceURL + "/api/post/" + id)
	if err == nil && resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		c.Data(http.StatusOK, "application/json", body)
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Content not found"})
}

func (h *GatewayHandler) GetYouTubeVideo(c *gin.Context) {
	id := c.Param("id")
	h.forwardToYouTube(c, "/api/video/"+id, nil)
}

func (h *GatewayHandler) GetInstagramPost(c *gin.Context) {
	id := c.Param("id")
	h.forwardToInstagram(c, "/api/post/"+id, nil)
}

func (h *GatewayHandler) forwardToYouTube(c *gin.Context, path string, body interface{}) {
	h.forwardRequest(c, h.youtubeServiceURL+path, body)
}

func (h *GatewayHandler) forwardToInstagram(c *gin.Context, path string, body interface{}) {
	h.forwardRequest(c, h.instagramServiceURL+path, body)
}

func (h *GatewayHandler) forwardRequest(c *gin.Context, url string, body interface{}) {
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal request body"})
			return
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(c.Request.Method, url, reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to service"})
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read response"})
		return
	}

	c.Data(resp.StatusCode, "application/json", respBody)
}