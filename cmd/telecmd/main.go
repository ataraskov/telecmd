package main

import (
	"log/slog"
	"os"

	"github.com/ataraskov/telecmd/internal/bot"
	"github.com/ataraskov/telecmd/internal/util"
)

func main() {
	slog.Info("Starting", "version", util.Version())

	token := os.Getenv("TOKEN")
	whitelist := os.Getenv("WHITELIST")
	verbose := len(os.Getenv("VERBOSE")) > 0

	cfg := bot.Config{
		Token:     token,
		Whitelist: util.ParseWhiteliest(whitelist),
		Verbose:   verbose,
	}

	b := bot.New(cfg)
	bot.Start(b)

	os.Exit(0)
}
