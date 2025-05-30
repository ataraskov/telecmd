package bot

import (
	"log/slog"
	"os"
	"time"

	"telecmd/internal/handler"

	tele "gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

// Config holds the bot configuration.
type Config struct {
	Token     string
	Whitelist []int64
	Verbose   bool
}

// New creates the Telegram bot.
func New(cfg Config) *tele.Bot {
	slog.Info("Initializing bot")

	pref := tele.Settings{
		Token:   cfg.Token,
		Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
		Verbose: cfg.Verbose,
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		slog.Error("Failed to create bot", "error", err)
		os.Exit(1)
	}

	b.Use(middleware.Recover())
	if len(cfg.Whitelist) > 0 {
		slog.Info("Adding whitelist", "whitelist", cfg.Whitelist)
		b.Use(middleware.Whitelist(cfg.Whitelist...))
	}

	b.Handle(tele.OnText, handler.TextHandler)

	return b
}

// Start starts the Telegram bot.
func Start(b *tele.Bot) {
	slog.Info("Bot starting...")
	b.Start()
}
