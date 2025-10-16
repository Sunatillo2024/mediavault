package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"youtube-service/internal/models"
	"youtube-service/internal/repository"

	"github.com/google/uuid"
)

type YouTubeService struct {
	repo *repository.VideoRepository
}

func NewYouTubeService(repo *repository.VideoRepository) *YouTubeService {
	return &YouTubeService{repo: repo}
}

func (s *YouTubeService) ExtractVideo(url string, downloadVideo, downloadSubtitles bool) (*models.Video, error) {
	// Use yt-dlp to get video info
	cmd := exec.Command("yt-dlp", "--dump-json", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to extract video info: %v", err)
	}

	var videoInfo map[string]interface{}
	if err := json.Unmarshal(output, &videoInfo); err != nil {
		return nil, fmt.Errorf("failed to parse video info: %v", err)
	}

	// Extract basic info
	videoID := videoInfo["id"].(string)
	title := videoInfo["title"].(string)
	description := videoInfo["description"].(string)
	channel := videoInfo["uploader"].(string)
	duration := int(videoInfo["duration"].(float64))

	// Get thumbnail URL
	thumbnail := videoInfo["thumbnail"].(string)

	// Extract available subtitles
	subtitles := s.extractSubtitles(videoInfo)

	// Create video model
	video := &models.Video{
		ID:           uuid.New().String(),
		YouTubeID:    videoID,
		Title:        title,
		Description:  description,
		ThumbnailURL: thumbnail,
		Duration:     duration,
		Channel:      channel,
		Subtitles:    subtitles,
	}

	// Download video if requested
	if downloadVideo {
		localPath, err := s.downloadVideo(url, video.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to download video: %v", err)
		}
		video.LocalPath = localPath
	}

	// Download subtitles if requested
	if downloadSubtitles {
		subtitlePaths, err := s.downloadSubtitles(url, video.ID, subtitles)
		if err != nil {
			return nil, fmt.Errorf("failed to download subtitles: %v", err)
		}
		video.Subtitles = subtitlePaths
	}

	// Save to database
	if err := s.repo.SaveVideo(video); err != nil {
		return nil, fmt.Errorf("failed to save video to database: %v", err)
	}

	return video, nil
}

func (s *YouTubeService) extractSubtitles(videoInfo map[string]interface{}) []string {
	subtitles := []string{}

	if subtitlesInfo, exists := videoInfo["subtitles"]; exists {
		if subtitleMap, ok := subtitlesInfo.(map[string]interface{}); ok {
			for lang := range subtitleMap {
				subtitles = append(subtitles, lang)
			}
		}
	}

	return subtitles
}

func (s *YouTubeService) downloadVideo(url, videoUUID string) (string, error) {
	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	videoDir := filepath.Join(storagePath, "youtube", "videos", videoUUID)
	if err := os.MkdirAll(videoDir, 0755); err != nil {
		return "", err
	}

	outputTemplate := filepath.Join(videoDir, "video.%(ext)s")
	cmd := exec.Command("yt-dlp", "-o", outputTemplate, url)

	if err := cmd.Run(); err != nil {
		return "", err
	}

	// Find the downloaded file
	files, err := os.ReadDir(videoDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			return filepath.Join(videoDir, file.Name()), nil
		}
	}

	return "", fmt.Errorf("no video file found after download")
}

func (s *YouTubeService) downloadSubtitles(url, videoUUID string, languages []string) ([]string, error) {
	if len(languages) == 0 {
		return []string{}, nil
	}

	storagePath := os.Getenv("STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./storage"
	}

	subtitleDir := filepath.Join(storagePath, "youtube", "subtitles", videoUUID)
	if err := os.MkdirAll(subtitleDir, 0755); err != nil {
		return nil, err
	}

	subtitlePaths := []string{}

	for _, lang := range languages {
		outputTemplate := filepath.Join(subtitleDir, fmt.Sprintf("subtitle.%s.vtt", lang))
		cmd := exec.Command("yt-dlp",
			"--write-sub",
			"--write-auto-sub",
			"--sub-lang", lang,
			"--skip-download",
			"-o", outputTemplate,
			url,
		)

		if err := cmd.Run(); err == nil {
			subtitlePaths = append(subtitlePaths, outputTemplate)
		}
	}

	return subtitlePaths, nil
}

func (s *YouTubeService) GetVideo(id string) (*models.Video, error) {
	return s.repo.GetVideoByID(id)
}

func (s *YouTubeService) GetSubtitles(id string) ([]string, error) {
	video, err := s.repo.GetVideoByID(id)
	if err != nil {
		return nil, err
	}
	return video.Subtitles, nil
}