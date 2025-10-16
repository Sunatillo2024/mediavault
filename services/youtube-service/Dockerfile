FROM golang:1.21-alpine

RUN apk add --no-cache \
    python3 \
    py3-pip \
    ffmpeg

RUN pip3 install yt-dlp

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main ./cmd/main.go

EXPOSE 8081

CMD ["./main"]