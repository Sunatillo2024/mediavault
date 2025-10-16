package models

import "time"

type Video struct {
	ID           string    `json:"id"`
	YouTubeID    string    `json:"youtube_id"`
	Title        string    `json:"title"`
	Description  string    `json:"description"`
	ThumbnailURL string    `json:"thumbnail_url"`
	Duration     int       `json:"duration"`
	Channel      string    `json:"channel"`
	Subtitles    []string  `json:"subtitles"`
	LocalPath    string    `json:"local_path,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}