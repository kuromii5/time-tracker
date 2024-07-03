package main

import (
	"log/slog"

	"github.com/kuromii5/time-tracker/internal/app"
	"github.com/kuromii5/time-tracker/internal/config"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

// @title Time Tracker
// @version 1.0
// @description This is time tracker app, where you can CRUD users, start/finish worklogs and watch them for users

// @host localhost:8080
// @BasePath /
func main() {
	cfg := config.MustLoad()
	logger := l.New(cfg.Env)

	// Create and configure the app
	application := app.New(
		logger,
		cfg.DbUrl,
		cfg.Port,
		cfg.RequestTimeout,
		cfg.IdleTimeout,
		cfg.ExternalAPIPort,
	)

	logger.Info("starting server", slog.Int("port", cfg.Port))

	// Start the app and handle graceful shutdown
	if err := application.Run(); err != nil {
		logger.Error("server failed", l.Err(err))
	}
}
