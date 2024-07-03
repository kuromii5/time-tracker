// Users handles user-related operations.
//
// swagger:tags
package user

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kuromii5/time-tracker/internal/models"
	"github.com/kuromii5/time-tracker/internal/utils"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type UsersGetter interface {
	Users(ctx context.Context, filter models.FilterBy, settings models.Pagination) ([]models.User, error)
}

type UsersResponse struct {
	Users []models.User `json:"users"`
}

// Render is used by chi/render to render the response.
func (ur UsersResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// Users handles retrieving a list of users.
// @Summary Get a list of users
// @Description Retrieve a list of users with optional filtering and pagination
// @Tags users
// @Accept json
// @Produce json
// @Param name query string false "Name"
// @Param surname query string false "Surname"
// @Param patronymic query string false "Patronymic"
// @Param address query string false "Address"
// @Param serie query string false "Passport Serie"
// @Param number query string false "Passport Number"
// @Param created_after query string false "Created After (timestamp)"
// @Param created_before query string false "Created Before (timestamp)"
// @Param limit query int false "Limit"
// @Param offset query int false "Offset"
// @Success 200 {object} UsersResponse "Successfully retrieved users"
// @Failure 500 {object} httperr.ErrResponse "Failed to get users"
// @Router /users [get]
func Users(logger *slog.Logger, usersGetter UsersGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "Users"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Parse query parameters for filtering
		filter := models.FilterBy{
			Name:           r.URL.Query().Get("name"),
			Surname:        r.URL.Query().Get("surname"),
			Patronymic:     r.URL.Query().Get("patronymic"),
			Address:        r.URL.Query().Get("address"),
			PassportSerie:  r.URL.Query().Get("serie"),
			PassportNumber: r.URL.Query().Get("number"),
			CreatedAfter:   utils.ParseQueryParamTime(r, "created_after"),
			CreatedBefore:  utils.ParseQueryParamTime(r, "created_before"),
		}

		// Parse query parameters for pagination
		pagination := models.Pagination{
			Limit:  utils.ParseQueryParamInt(r, "limit"),
			Offset: utils.ParseQueryParamInt(r, "offset"),
		}

		log.Debug("received request",
			slog.Any("filter", filter),
			slog.Any("pagination", pagination),
		)

		// Get users from the database
		users, err := usersGetter.Users(r.Context(), filter, pagination)
		if err != nil {
			log.Error("failed to get users", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		log.Info("fetched users", slog.Int("count", len(users)))

		// Write response
		resp := UsersResponse{Users: users}
		if err := render.Render(w, r, resp); err != nil {
			log.Error("failed to render response", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
		}
	}
}
