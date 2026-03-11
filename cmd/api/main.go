package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/MotiurRahmanSany/student-management-api/internal/config"
)

func setupLogger() {
	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo, // Only Info and above (no Debug noise in production)
	}

	if os.Getenv("APP_ENV") == "production" {
		// JSON format — for log aggregation tools (Datadog, Loki, CloudWatch, etc.)
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// Human-readable text format — for local development
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	// SetDefault makes this handler the global default
	// After this, all slog.Info/Warn/Error calls anywhere in the program use it
	slog.SetDefault(slog.New(handler))
}

func main() {

	// setup logger before any other logging calls
	setupLogger()

	fmt.Println("Loading configuration...")
	config := config.GetConfig()
	fmt.Printf("Starting %s version %s on port %d\n", config.AppName, config.Version, config.HttpPort)
	fmt.Printf("Database host: %s, port: %d, user: %s\n", config.Db.Host, config.Db.Port, config.Db.User)

	serve(config)
}
