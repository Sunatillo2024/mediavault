package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"youtube-service/internal/models"

	_ "github.com/lib/pq"
)

type VideoRepository struct {
	db *sql.DB
}

func NewPostgresDB() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create table if not exists
	if err := createYouTubeTable(db); err != nil {
		return nil, err
	}

	return db, nil
}

func createYouTubeTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS youtube_videos (
		id VARCHAR(255) PRIMARY KEY,
		youtube_id VARCHAR(255) NOT NULL,
		title TEXT,
		description TEXT,
		thumbnail_url TEXT,
		duration INTEGER,
		channel VARCHAR(255),
		subtitles JSONB,
		local_path TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.Exec(createTableSQL)
	return err
}

func NewVideoRepository(db *sql.DB) *VideoRepository {
	return &VideoRepository{db: db}
}

func (r *VideoRepository) SaveVideo(video *models.Video) error {
	subtitlesJSON, err := json.Marshal(video.Subtitles)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO youtube_videos (id, youtube_id, title, description, thumbnail_url, duration, channel, subtitles, local_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err = r.db.Exec(query,
		video.ID,
		video.YouTubeID,
		video.Title,
		video.Description,
		video.ThumbnailURL,
		video.Duration,
		video.Channel,
		subtitlesJSON,
		video.LocalPath,
	)

	return err
}

func (r *VideoRepository) GetVideoByID(id string) (*models.Video, error) {
	query := `
		SELECT id, youtube_id, title, description, thumbnail_url, duration, channel, subtitles, local_path, created_at
		FROM youtube_videos
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	var video models.Video
	var subtitlesJSON string
	var createdAt sql.NullTime

	err := row.Scan(
		&video.ID,
		&video.YouTubeID,
		&video.Title,
		&video.Description,
		&video.ThumbnailURL,
		&video.Duration,
		&video.Channel,
		&subtitlesJSON,
		&video.LocalPath,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video not found")
		}
		return nil, err
	}

	// Parse subtitles JSON
	if err := json.Unmarshal([]byte(subtitlesJSON), &video.Subtitles); err != nil {
		log.Printf("Warning: Failed to parse subtitles JSON: %v", err)
		video.Subtitles = []string{}
	}

	if createdAt.Valid {
		video.CreatedAt = createdAt.Time
	}

	return &video, nil
}

func (r *VideoRepository) GetVideoByYouTubeID(youtubeID string) (*models.Video, error) {
	query := `
		SELECT id, youtube_id, title, description, thumbnail_url, duration, channel, subtitles, local_path, created_at
		FROM youtube_videos
		WHERE youtube_id = $1
	`

	row := r.db.QueryRow(query, youtubeID)

	var video models.Video
	var subtitlesJSON string
	var createdAt sql.NullTime

	err := row.Scan(
		&video.ID,
		&video.YouTubeID,
		&video.Title,
		&video.Description,
		&video.ThumbnailURL,
		&video.Duration,
		&video.Channel,
		&subtitlesJSON,
		&video.LocalPath,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("video not found")
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(subtitlesJSON), &video.Subtitles); err != nil {
		log.Printf("Warning: Failed to parse subtitles JSON: %v", err)
		video.Subtitles = []string{}
	}

	if createdAt.Valid {
		video.CreatedAt = createdAt.Time
	}

	return &video, nil
}package repository

