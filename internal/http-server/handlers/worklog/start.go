package worklog

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type WorklogStarter interface {
	StartWorklog(ctx context.Context, task string, userID int32) (int32, error)
}

type StartWorklogRequest struct {
	Task   string `json:"task"`
	UserID int32  `json:"user_id"`
}

type StartWorklogResponse struct {
	WorklogID int32 `json:"worklog_id"`
}

// @Summary Start a worklog
// @Description Start a new worklog for a specified user with a given task
// @Tags worklogs
// @Accept json
// @Produce json
// @Param request body StartWorklogRequest true "Start Worklog Request"
// @Success 201 {object} StartWorklogResponse "Successfully started worklog"
// @Failure 400 {object} httperr.ErrResponse "Invalid request payload"
// @Failure 500 {object} httperr.ErrResponse "Internal server error"
// @Router /worklogs/start [post]
func StartWorklog(logger *slog.Logger, worklogStarter WorklogStarter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "StartWorklog"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req StartWorklogRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				render.Render(w, r, httperr.ErrInvalidRequest(err))
				return
			}
			log.Error("failed to decode request body", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}
		defer r.Body.Close()

		// create record in DB
		worklogID, err := worklogStarter.StartWorklog(r.Context(), req.Task, req.UserID)
		if err != nil {
			log.Error("failed to start worklog", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		resp := StartWorklogResponse{WorklogID: worklogID}
		log.Info("worklog started successfully", slog.Int("worklog_id", int(worklogID)))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, resp)
	}
}
