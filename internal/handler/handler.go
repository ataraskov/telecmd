package handler

import (
	"log/slog"
	"strings"

	"github.com/ataraskov/telecmd/internal/features/downloader"

	tele "gopkg.in/telebot.v4"
)

// TextHandler handles incoming text messages.
func TextHandler(c tele.Context) error {
	logger := slog.With("command", "text", "text", c.Text())

	switch {
	case strings.HasPrefix(c.Text(), "http"):
		logger.Info("Detected a link")
		return DownloadHandler(c)
	default:
		logger.Info("Received a text message")
		return c.Send("Give me a link ;)")
	}
}

// DownloadHandler handles the /download command.
func DownloadHandler(c tele.Context) error {
	logger := slog.With("command", "download")

	url := c.Text()
	if url == "" {
		return c.Send("Please provide a URL. Usage: <url>")
	}

	c.Send("Downloading...")
	d := downloader.New()
	filename, _, err := d.DownloadVideo(url)
	if err != nil || filename == "" {
		logger.Error("Error downloading video", "url", url, "filename", filename, "error", err)
		return c.Send("Error downloading video: " + err.Error())
	}
	defer d.RemoveFile(filename)

	logger.Info("Downloaded file", "filename", filename, "url", url)
	video := &tele.Video{File: tele.FromDisk(filename)}
	err = c.Send(video)
	if err != nil {
		logger.Error("Error sending video", "filename", filename, "error", err)
		return c.Send("Error sending video: " + err.Error())
	}

	return nil
}
