package util

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ParseWhiteliest parses a comma-separated whitelist to a slice of int64.
func ParseWhiteliest(whitelist string) []int64 {
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

// RemoveFile removes the file after sending
func RemoveFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		slog.Error("Error removing file", "filename", filename, "error", err)
	} else {
		slog.Info("File removed successfully", "filename", filename)
	}
}

// RunCommand executes an external command with given arguments.
// It returns the standard output, standard error, and any execution error.
func RunCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var cmdout bytes.Buffer
	cmd.Stdout = &cmdout
	cmd.Stderr = &cmdout

	err := cmd.Run()

	return cmdout.String(), err
}

// DownloadVideo downloads a video using yt-dlp
func DownloadVideo(url string) (string, string, error) {
	// Generate a temporary filename to avoid conflicts
	tmpfile, err := os.CreateTemp("", "telecmd_")
	if err != nil {
		return "", "", err
	}

	filename := tmpfile.Name() + ".mp4"
	output, err := RunCommand(
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
