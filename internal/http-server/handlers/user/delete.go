// Users handles user-related operations.
//
// swagger:tags
package user

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type UserDeleter interface {
	DeleteUser(ctx context.Context, id int32) error
}

// DeleteUser handles the deletion of a user.
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} httperr.ErrResponse "Invalid user ID"
// @Failure 500 {object} httperr.ErrResponse "Failed to delete user"
// @Router /users/{id} [delete]
func DeleteUser(logger *slog.Logger, userDeleter UserDeleter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "DeleteUser"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Parse user ID from the URL
		idStr := chi.URLParam(r, "id")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid user ID", slog.String("user_id", idStr), l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(errors.New("invalid user ID")))
			return
		}

		// Delete user from the database
		if err := userDeleter.DeleteUser(r.Context(), int32(userId)); err != nil {
			log.Error("failed to delete user", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		log.Info("deleted user", slog.Int("user_id", int(userId)))

		w.WriteHeader(http.StatusNoContent)
	}
}
