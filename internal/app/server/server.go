package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/kuromii5/time-tracker/docs"
	"github.com/kuromii5/time-tracker/internal/http-server/handlers/user"
	"github.com/kuromii5/time-tracker/internal/http-server/handlers/worklog"
	mwlog "github.com/kuromii5/time-tracker/internal/http-server/middleware/mw_log"
	"github.com/kuromii5/time-tracker/internal/repo"
	httpSwagger "github.com/swaggo/http-swagger"
)

func New(
	logger *slog.Logger,
	port int,
	reqTimeout, idleTimeout time.Duration,
	db *repo.DB,
	externalAPIPort int,
) *http.Server {
	r := chi.NewRouter()

	applyMiddlewares(r, logger)
	setupRoutes(r, logger, db, externalAPIPort)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  reqTimeout,
		WriteTimeout: reqTimeout,
		IdleTimeout:  idleTimeout,
	}

	return srv
}

func applyMiddlewares(r *chi.Mux, logger *slog.Logger) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(mwlog.New(logger)) // use custom logger for http requests
	r.Use(middleware.Recoverer)
}

func setupRoutes(r *chi.Mux, logger *slog.Logger, db *repo.DB, extAPIPort int) {
	// use swagger
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), // The url pointing to API definition
	))

	// user routes
	r.Get("/users", user.Users(logger, db))
	r.Post("/users", user.CreateUser(logger, db, extAPIPort))
	r.Patch("/users/{id}", user.UpdateUser(logger, db))
	r.Delete("/users/{id}", user.DeleteUser(logger, db))

	// worklog routes
	r.Get("/users/{userID}/worklogs", worklog.Worklogs(logger, db))
	r.Post("/worklogs/start", worklog.StartWorklog(logger, db))
	r.Patch("/worklogs/finish/{id}", worklog.FinishWorklog(logger, db))
}
