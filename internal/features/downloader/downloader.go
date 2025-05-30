package downloader

import (
	"bytes"
	"log/slog"
	"os"
	"os/exec"
)

type Downloader struct {
	Format string
}

func New() *Downloader {
	defaultFormat := "mp4,res:720,filesize<50M"
	return &Downloader{
		Format: defaultFormat,
	}
}

func (d *Downloader) WithFormat(format string) *Downloader {
	return &Downloader{
		Format: format,
	}
}

// RemoveFile removes the file after sending
func (d *Downloader) RemoveFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		slog.Error("Error removing file", "filename", filename, "error", err)
	} else {
		slog.Info("File removed successfully", "filename", filename)
	}
}

// DownloadVideo downloads a video using yt-dlp
// returns the filename, output, and any error
func (d *Downloader) DownloadVideo(url string) (string, string, error) {
	// Generate a temporary filename to avoid conflicts
	tmpfile, err := os.CreateTemp("", "telecmd_")
	if err != nil {
		return "", "", err
	}

	filename := tmpfile.Name() + ".mp4"
	output, err := d.runCommand(
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

// runCommand executes an external command with given arguments.
// It returns the standard output, standard error, and any execution error.
func (d *Downloader) runCommand(name string, arg ...string) (string, error) {
	cmd := exec.Command(name, arg...)
	var cmdout bytes.Buffer
	cmd.Stdout = &cmdout
	cmd.Stderr = &cmdout

	err := cmd.Run()

	return cmdout.String(), err
}
