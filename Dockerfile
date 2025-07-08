# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o telecmd cmd/telecmd/main.go

# Final stage
FROM python:3.13-alpine

# Install yt-dlp
RUN wget -O /usr/local/bin/yt-dlp https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp && \
    chmod +x /usr/local/bin/yt-dlp

# Copy scripts
COPY scripts/update-yt-dlp.sh /usr/local/bin/update-yt-dlp.sh
COPY scripts/start.sh         /usr/local/bin/start.sh
RUN chmod +x /usr/local/bin/update-yt-dlp.sh /usr/local/bin/start.sh
RUN echo "@daily /usr/local/bin/update-yt-dlp.sh" > /etc/crontabs/root

# Copy the binary from builder stage
COPY --from=builder /app/telecmd /usr/local/bin/telecmd

WORKDIR /app

CMD ["/usr/local/bin/start.sh"]
