package util

import (
	"log/slog"
	"runtime/debug"
	"strconv"
	"strings"
)

func Version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "unknown version (no build info)"
	}

	version := info.Main.Version
	if version == "" || version == "(devel)" {
		// Try to find VCS info if main version is missing or devel
		revision := ""
		buildTime := ""
		for _, setting := range info.Settings {
			switch setting.Key {
			case "vcs.revision":
				revision = setting.Value
			case "vcs.time":
				buildTime = setting.Value
			}
		}
		if revision != "" {
			version = "devel (" + revision
			if buildTime != "" {
				version += " @ " + buildTime
			}
			version += ")"
		} else {
			// Fallback if no version and no VCS info
			version = "unknown version (devel)"
		}
	}

	return version
}

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
