# MediaVault

A microservices-based platform for extracting and storing content from YouTube and Instagram. MediaVault provides a unified API gateway for downloading videos, metadata, and subtitles from multiple social media platforms.

## 📋 Table of Contents

- [Features](#features)
- [Architecture](#architecture)
- [Tech Stack](#tech-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [API Documentation](#api-documentation)
- [Testing](#testing)
- [Project Structure](#project-structure)
- [Contributing](#contributing)
- [License](#license)

## ✨ Features

- **Multi-Platform Support**: Extract content from YouTube and Instagram
- **Microservices Architecture**: Independent, scalable services
- **Unified API Gateway**: Single entry point for all services
- **Metadata Extraction**: Retrieve titles, descriptions, hashtags, and more
- **Media Download**: Download videos, images, and subtitles
- **PostgreSQL Storage**: Persistent storage for all extracted content
- **Redis Caching**: Fast access to frequently requested data
- **Docker Support**: Easy deployment with Docker Compose
- **Health Monitoring**: Service health checks and status monitoring

## 🏗 Architecture

```
┌─────────────┐
│   Client    │
└──────┬──────┘
       │
       ▼
┌─────────────────┐
│  API Gateway    │ (Port 8080)
│    (Go/Gin)     │
└────────┬────────┘
         │
    ┌────┴────┐
    ▼         ▼
┌─────────┐ ┌──────────┐
│ YouTube │ │Instagram │
│ Service │ │ Service  │
│(Go/Gin) │ │(Python)  │
│Port 8081│ │Port 8082 │
└────┬────┘ └────┬─────┘
     │           │
     └─────┬─────┘
           ▼
    ┌────────────┐
    │ PostgreSQL │
    │   Redis    │
    └────────────┘
```

## 🛠 Tech Stack

### Backend Services

- **Gateway Service**: Go 1.21, Gin Web Framework
- **YouTube Service**: Go 1.21, Gin, yt-dlp
- **Instagram Service**: Python 3.11, FastAPI, Instaloader

### Infrastructure

- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Containerization**: Docker, Docker Compose
- **API Framework**: Gin (Go), FastAPI (Python)

### External Tools

- **yt-dlp**: YouTube video extraction
- **Instaloader**: Instagram content extraction
- **FFmpeg**: Media processing

## 📦 Prerequisites

Before you begin, ensure you have the following installed:

- **Docker** (v20.10+)
- **Docker Compose** (v2.0+)
- **Git**

Optional for local development:
- Go 1.21+
- Python 3.11+
- PostgreSQL 15+
- Redis 7+

## 🚀 Installation

### Quick Start with Docker

1. **Clone the repository**
```bash
git clone https://github.com/yourusername/mediavault.git
cd mediavault
```

2. **Create environment file**
```bash
cp .env.example .env
```

3. **Update environment variables** (optional)
Edit `.env` file with your preferred configurations:
```bash
# Gateway
GATEWAY_PORT=8080
JWT_SECRET=your_secure_secret_key

# Database
DB_PASSWORD=your_secure_password

# Services
STORAGE_PATH=/app/storage
```

4. **Build and start services**
```bash
docker-compose up --build
```

5. **Verify services are running**
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "services": {
    "youtube": "ok",
    "instagram": "ok"
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### Local Development Setup

#### Gateway Service

```bash
cd gateway
go mod download
go run cmd/main.go
```

#### YouTube Service

```bash
cd services/youtube-service
go mod download
go run cmd/main.go
```

#### Instagram Service

```bash
cd services/instagram-service
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8082
```

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `GATEWAY_PORT` | Gateway service port | `8080` |
| `YOUTUBE_SERVICE_URL` | YouTube service URL | `http://youtube-service:8081` |
| `INSTAGRAM_SERVICE_URL` | Instagram service URL | `http://instagram-service:8082` |
| `JWT_SECRET` | JWT secret key | `your_secret_key` |
| `DB_HOST` | PostgreSQL host | `postgres` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_NAME` | Database name | `mediavault` |
| `DB_USER` | Database user | `admin` |
| `DB_PASSWORD` | Database password | `password123` |
| `REDIS_HOST` | Redis host | `redis` |
| `REDIS_PORT` | Redis port | `6379` |
| `STORAGE_PATH` | Media storage path | `/app/storage` |

### Database Schema

The database schema is automatically created on service startup:

- **youtube_videos**: Stores YouTube video metadata
- **instagram_posts**: Stores Instagram post metadata

## 📖 Usage

### Extract YouTube Video

```bash
curl -X POST http://localhost:8080/api/youtube/extract \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
    "download_video": true,
    "download_subtitles": true
  }'
```

### Extract Instagram Post

```bash
curl -X POST http://localhost:8080/api/instagram/extract \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.instagram.com/p/ABC123/",
    "download_media": true
  }'
```

### Get Content by ID

```bash
curl http://localhost:8080/api/content/{id}
```

### Platform-Specific Endpoints

**YouTube:**
```bash
# Get video details
curl http://localhost:8080/api/youtube/video/{id}

# Get subtitles
curl http://localhost:8080/api/youtube/subtitles/{id}
```

**Instagram:**
```bash
# Get post details
curl http://localhost:8080/api/instagram/post/{id}

# Get media paths
curl http://localhost:8080/api/instagram/media/{id}
```

## 📚 API Documentation

### Gateway Endpoints

#### Health Check
```http
GET /health
```
Returns the health status of all services.

#### Universal Extract
```http
POST /api/extract
Content-Type: application/json

{
  "url": "string",
  "platform": "youtube|instagram",
  "download_video": boolean,
  "download_subtitles": boolean,
  "download_media": boolean
}
```

#### Get Content
```http
GET /api/content/:id
```
Retrieves content from any platform by ID.

### YouTube Service Endpoints

#### Extract Video
```http
POST /api/extract
Content-Type: application/json

{
  "url": "string",
  "download_video": boolean,
  "download_subtitles": boolean
}
```

#### Get Video
```http
GET /api/video/:id
```

#### Get Subtitles
```http
GET /api/subtitles/:id
```

### Instagram Service Endpoints

#### Extract Post
```http
POST /api/extract
Content-Type: application/json

{
  "url": "string",
  "download_media": boolean
}
```

#### Get Post
```http
GET /api/post/:id
```

#### Get Media Paths
```http
GET /api/media/:id
```

## 🧪 Testing

### Manual Testing

1. **Test Gateway Health**
```bash
curl http://localhost:8080/health
```

2. **Test YouTube Extraction**
```bash
curl -X POST http://localhost:8080/api/youtube/extract \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.youtube.com/watch?v=jNQXAC9IVRw",
    "download_video": false,
    "download_subtitles": false
  }'
```

3. **Test Instagram Extraction**
```bash
curl -X POST http://localhost:8080/api/instagram/extract \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://www.instagram.com/p/C1234567890/",
    "download_media": false
  }'
```

### Automated Testing

Run unit tests for each service:

**Gateway:**
```bash
cd gateway
go test ./...
```

**YouTube Service:**
```bash
cd services/youtube-service
go test ./...
```

**Instagram Service:**
```bash
cd services/instagram-service
pytest
```

### Integration Testing

```bash
# Start all services
docker-compose up -d

# Run integration tests
./scripts/integration-tests.sh
```

## 📁 Project Structure

```
mediavault/
├── docker-compose.yml
├── .env.example
├── init.sql
├── README.md
├── gateway/
│   ├── Dockerfile
│   ├── go.mod
│   ├── go.sum
│   ├── cmd/
│   │   └── main.go
│   └── internal/
│       └── handler/
│           └── gateway_handler.go
├── services/
│   ├── youtube-service/
│   │   ├── Dockerfile
│   │   ├── go.mod
│   │   ├── go.sum
│   │   ├── cmd/
│   │   │   └── main.go
│   │   └── internal/
│   │       ├── handler/
│   │       ├── models/
│   │       ├── repository/
│   │       └── service/
│   └── instagram-service/
│       ├── Dockerfile
│       ├── requirements.txt
│       └── app/
│           ├── main.py
│           ├── api/
│           ├── models/
│           ├── services/
│           └── database.py
└── storage/
    ├── youtube/
    │   ├── videos/
    │   └── subtitles/
    └── instagram/
        └── media/
```

## 🔧 Troubleshooting

### Services won't start

1. Check Docker logs:
```bash
docker-compose logs gateway
docker-compose logs youtube-service
docker-compose logs instagram-service
```

2. Verify ports are not in use:
```bash
lsof -i :8080
lsof -i :8081
lsof -i :8082
```

3. Rebuild containers:
```bash
docker-compose down -v
docker-compose up --build
```

### Database connection issues

```bash
# Check PostgreSQL logs
docker-compose logs postgres

# Connect to database manually
docker-compose exec postgres psql -U admin -d mediavault
```

### YouTube extraction fails

- Ensure yt-dlp is up to date (rebuilt Docker image)
- Check if the video is available in your region
- Verify the video URL is correct

### Instagram extraction fails

- Instagram may rate-limit requests
- Some private accounts require authentication
- Check Instaloader logs for specific errors

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Guidelines

- Follow Go and Python best practices
- Write unit tests for new features
- Update documentation as needed
- Use meaningful commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 👥 Authors

- Your Name - Initial work

## 🙏 Acknowledgments

- [yt-dlp](https://github.com/yt-dlp/yt-dlp) - YouTube extractor
- [Instaloader](https://instaloader.github.io/) - Instagram extractor
- [Gin](https://gin-gonic.com/) - Go web framework
- [FastAPI](https://fastapi.tiangolo.com/) - Python web framework

## 📞 Support

For support, email support@mediavault.com or open an issue on GitHub.

## 🗺 Roadmap

- [ ] Add authentication and user management
- [ ] Support for additional platforms (TikTok, Twitter)
- [ ] Implement rate limiting
- [ ] Add webhook notifications
- [ ] Create web dashboard
- [ ] Add video transcoding
- [ ] Implement search functionality
- [ ] Add bulk extraction support

---

Made with ❤️ by the MediaVault Team
