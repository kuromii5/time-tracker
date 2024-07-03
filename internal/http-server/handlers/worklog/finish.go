package worklog

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kuromii5/time-tracker/internal/repo"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type WorklogFinisher interface {
	FinishWorklog(ctx context.Context, worklogID int32) error
}

// @Summary Finish a worklog
// @Description Finish a worklog with the specified ID
// @Tags worklogs
// @Accept json
// @Produce json
// @Param id path int true "Worklog ID"
// @Success 204 "No Content"
// @Failure 400 {object} httperr.ErrResponse "Invalid worklog ID"
// @Failure 500 {object} httperr.ErrResponse "Failed to finish worklog"
// @Router /worklogs/finish/{id} [patch]
func FinishWorklog(logger *slog.Logger, worklogFinisher WorklogFinisher) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "FinishWorklog"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// update record in DB
		worklogID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("failed to parse worklog ID", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		err = worklogFinisher.FinishWorklog(r.Context(), int32(worklogID))
		if err != nil {
			if err == repo.ErrAlreadyDone {
				log.Warn("this worklog was already finished", l.Err(err))

				render.Render(w, r, httperr.ErrConflict(err))
				return
			}
			log.Error("failed to finish worklog", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		log.Info("worklog finished successfully", slog.Int("worklog_id", worklogID))

		w.WriteHeader(http.StatusNoContent)
	}
}
