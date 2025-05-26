package main

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func main() {
	slog.Info("Starting")
	if len(os.Getenv("TOKEN")) == 0 {
		slog.Warn("Environment variable TOKEN is not set.")
	}
	token := os.Getenv("TOKEN")
	whitelist := parseWhiteliest(os.Getenv("WHITELIST"))
	verbose := len(os.Getenv("VERBOSE")) > 0

	pref := tele.Settings{
		Token:   token,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: verbose,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("Failed to create bot", "error", err)
		os.Exit(1)
	}

	b.Use(middleware.Recover())
	if len(whitelist) > 0 {
		slog.Info("Adding whitelist", "whitelist", whitelist)
		b.Use(middleware.Whitelist(whitelist...))
	}

	b.Handle(tele.OnText, textHandler)

	slog.Info("Bot starting...")
	b.Start()
}

// parseWhiteliest parses a comma-separated whitelist to a slice of int64.
func parseWhiteliest(whitelist string) []int64 {
	parts := strings.Split(whitelist, ",")
	result := make([]int64, 0, len(parts))

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil {
			slog.Warn("Invalid whitelist entry", "entry", part, "error", err)
			continue
		}

		result = append(result, id)
	}

	return result
}

// textHandler handles incoming text messages.
func textHandler(c tele.Context) error {
	logger := slog.With("command", "text", "text", c.Text())

	switch {
	case strings.HasPrefix(c.Text(), "http"):
		logger.Info("Detected a link")
		return downloadHandler(c)
	default:
		logger.Info("Received a text message")
		return c.Send("Give me a link ;)")
	}
}

// downloadHandler handles the /download command.
func downloadHandler(c tele.Context) error {
	logger := slog.With("command", "download")

	url := c.Text()
	if url == "" {
		return c.Send("Please provide a URL. Usage: <url>")
	}

	c.Send("Downloading...")
	filename, _, err := downloadVideo(url)
	if err != nil || filename == "" {
		logger.Error("Error downloading video", "url", url, "filename", filename, "error", err)
		return c.Send("Error downloading video: " + err.Error())
	}
	defer removeFile(filename)

	logger.Info("Downloaded file", "filename", filename, "url", url)
	video := &tele.Video{File: tele.FromDisk(filename)}
	err = c.Send(video)
	if err != nil {
		logger.Error("Error sending video", "filename", filename, "error", err)
		return c.Send("Error sending video: " + err.Error())
	}

	return nil
}

func removeFile(filename string) {
	// Remove the file after sending
	err := os.Remove(filename)
	if err != nil {
		slog.Error("Error removing file", "filename", filename, "error", err)
	} else {
		slog.Info("File removed successfully", "filename", filename)
	}
}

// runCommand executes an external command with given arguments.
// It returns the standard output, standard error, and any execution error.
func runCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var cmdout bytes.Buffer
	cmd.Stdout = &cmdout
	cmd.Stderr = &cmdout

	err := cmd.Run()

	return cmdout.String(), err
}

// downloadVideo downloads a video using yt-dlp
func downloadVideo(url string) (string, string, error) {
	// Generate a temporary filename to avoid conflicts
	tmpfile, err := os.CreateTemp("", "telecmd_")
	if err != nil {
		return "", "", err
	}

	filename := tmpfile.Name() + ".mp4"
	output, err := runCommand(
		"yt-dlp",
		"--format", "mp4,res:720,filesize<50M",
		"--output", filename,
		url,
	)
	if err != nil {
		return "", "", err
	}

	return filename, output, nil
}
