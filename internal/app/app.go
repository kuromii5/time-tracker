package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kuromii5/time-tracker/internal/app/server"
	"github.com/kuromii5/time-tracker/internal/repo"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type App struct {
	logger *slog.Logger
	server *http.Server
	db     *repo.DB
}

func New(
	logger *slog.Logger,
	dbUrl string,
	port int,
	reqTimeout, idleTimeout time.Duration,
	externalAPIPort int,
) *App {
	db, err := repo.New(dbUrl, logger)
	if err != nil {
		log.Fatalf("Failed to connect to db: %v", err)
	}

	server := server.New(logger, port, reqTimeout, idleTimeout, db, externalAPIPort)

	return &App{
		logger: logger,
		server: server,
		db:     db,
	}
}

func (a *App) Run() error {
	go func() {
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error("server failed", l.Err(err))
		}
	}()

	// Set up graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	a.logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		a.logger.Error("server shutdown error", l.Err(err))

		return err
	}

	a.logger.Info("server stopped gracefully")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	// Close db first
	a.db.Close()
	return a.server.Shutdown(ctx)
}
