package worklog

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kuromii5/time-tracker/internal/models"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type WorklogsGetter interface {
	Worklogs(ctx context.Context, userID int32, startDate, endDate time.Time) ([]models.Worklog, error)
}

type WorklogsRequest struct {
	StartDate time.Time `json:"start_date" description:"Start date in the format YYYY-MM-DDTHH:MM:SSZ (ISO 8601)"`
	EndDate   time.Time `json:"end_date" description:"End date in the format YYYY-MM-DDTHH:MM:SSZ (ISO 8601)"`
}

type WorklogResponse struct {
	ID        int32  `json:"id"`
	UserID    int32  `json:"user_id"`
	Task      string `json:"task"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
	Duration  string `json:"duration"`
}

func formatDuration(duration time.Duration) string {
	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}
func formatTime(t time.Time) string {
	return t.Format("2006-01-02, 15:04:05")
}

// @Summary Get worklogs for a user
// @Description Get worklogs for a user within a specified date range. Time format should be YYYY-MM-DDTHH:MM:SSZ (ISO 8601).
// @Tags worklogs
// @Accept json
// @Produce json
// @Param userID path int true "User ID"
// @Param request body WorklogsRequest true "Worklogs Request"
// @Success 200 {array} WorklogResponse "List of worklogs"
// @Failure 400 {object} httperr.ErrResponse "Invalid request payload or user ID"
// @Failure 500 {object} httperr.ErrResponse "Failed to get worklogs"
// @Router /users/{userID}/worklogs [get]
func Worklogs(logger *slog.Logger, worklogsGetter WorklogsGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "Worklogs"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		userID, err := strconv.Atoi(chi.URLParam(r, "userID"))
		if err != nil {
			log.Error("invalid user ID", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		var req WorklogsRequest
		err = json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Error("failed to decode request body", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		if req.EndDate.Before(req.StartDate) {
			err := fmt.Errorf("end date should be after start date")
			log.Error("invalid date range", l.Err(err))
			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		log.Debug("start_date", slog.Time("start_date", req.StartDate))
		log.Debug("end_date", slog.Time("end_date", req.EndDate))

		worklogs, err := worklogsGetter.Worklogs(r.Context(), int32(userID), req.StartDate, req.EndDate)
		if err != nil {
			log.Error("failed to get worklogs", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		var resp []WorklogResponse
		for _, wl := range worklogs {
			duration := formatDuration(wl.FinishedAt.Sub(wl.StartedAt))
			startTime := formatTime(wl.StartedAt)
			endTime := formatTime(wl.FinishedAt)

			wr := WorklogResponse{
				ID:        wl.ID,
				UserID:    wl.UserID,
				Task:      wl.Task,
				StartTime: startTime,
				EndTime:   endTime,
				Duration:  duration,
			}
			resp = append(resp, wr)
		}

		render.JSON(w, r, resp)

		log.Info("worklogs retrieved successfully", slog.Int("count", len(worklogs)))
	}
}
