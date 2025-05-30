package util

import (
	"fmt"
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

	if info.Main.Version == "" {
		return "devel"
	}

	return fmt.Sprintf("%s (%s)",
		info.Main.Version,
		info.Main.Sum,
	)
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
