// Users handles user-related operations.
//
// swagger:tags
package user

import (
	"context"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kuromii5/time-tracker/internal/models"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type UserUpdater interface {
	UpdateUser(ctx context.Context, user models.User) error
}

type UpdateUserRequest struct {
	Passport struct {
		Serie  string `json:"serie"`
		Number string `json:"number"`
	} `json:"passport"`
	People struct {
		Name       string `json:"name"`
		Surname    string `json:"surname"`
		Patronymic string `json:"patronymic"`
		Address    string `json:"address"`
	} `json:"people"`
}

// UpdateUser handles updating an existing user.
// @Summary Update an existing user
// @Description Update a user's details using the provided information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update User Request"
// @Success 204 "Successfully updated user"
// @Failure 400 {object} httperr.ErrResponse "Invalid request payload"
// @Failure 500 {object} httperr.ErrResponse "Internal server error"
// @Router /users/{id} [patch]
func UpdateUser(logger *slog.Logger, userUpdater UserUpdater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "UpdateUser"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Parse the request body into an UpdateUserRequest object
		var req UpdateUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			log.Error("failed to decode request body", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		// Parse user ID from the URL
		idStr := chi.URLParam(r, "id")
		userId, err := strconv.Atoi(idStr)
		if err != nil {
			log.Error("invalid user ID", slog.String("user_id", idStr), l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		// Prepare the user object for update
		user := models.User{
			ID: int32(userId),
			Passport: models.Passport{
				Serie:  req.Passport.Serie,
				Number: req.Passport.Number,
			},
			People: models.People{
				Name:       req.People.Name,
				Surname:    req.People.Surname,
				Patronymic: req.People.Patronymic,
				Address:    req.People.Address,
			},
		}

		// Update user in the database
		if err := userUpdater.UpdateUser(r.Context(), user); err != nil {
			log.Error("failed to update user", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		log.Info("updated user", slog.Int("user_id", int(user.ID)))

		w.WriteHeader(http.StatusNoContent)
	}
}
