# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telecmd .

# Final stage
FROM python:3.13-alpine

# Install yt-dlp
RUN wget -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp && \
    chmod +x /usr/local/bin/yt-dlp

# Copy the binary from builder stage
COPY --from=builder /app/telecmd /usr/local/bin/telecmd

WORKDIR /app

CMD ["telecmd"]