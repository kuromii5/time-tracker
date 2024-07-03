// Users handles user-related operations.
//
// swagger:tags
package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/kuromii5/time-tracker/internal/models"
	"github.com/kuromii5/time-tracker/internal/utils"
	httperr "github.com/kuromii5/time-tracker/pkg/http-errors"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

type UserCreator interface {
	CreateUser(ctx context.Context, user models.User) (int32, error)
}

type CreateUserRequest struct {
	PassportNumber string `json:"passportNumber"`
}
type CreateUserResponse struct {
	UserID int32 `json:"user_id"`
}

// CreateUser handles the creation of a new user.
// @Summary Create a new user
// @Description Create a new user with provided passport number and fetch additional data from an external API
// @Tags users
// @Accept json
// @Produce json
// @Param extAPIPort query int true "External API Port" default(8081)
// @Param request body CreateUserRequest true "Create User Request"
// @Success 201 {object} CreateUserResponse "Successfully created user"
// @Failure 400 {object} httperr.ErrResponse "Invalid request payload"
// @Failure 500 {object} httperr.ErrResponse "Internal server error"
// @Router /users [post]
func CreateUser(logger *slog.Logger, userCreator UserCreator, extAPIPort int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log := logger.With(
			slog.String("handler", "CreateUser"),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req CreateUserRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")

				render.Render(w, r, httperr.ErrInvalidRequest(errors.New("request body is empty")))
				return
			}
			log.Error("failed to decode request body", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}
		defer r.Body.Close()

		log.Debug("passport data", slog.String("passportNumber", req.PassportNumber))

		passport, err := utils.ParsePassportData(req.PassportNumber)
		if err != nil {
			log.Error("failed to parse passport data", l.Err(err))

			render.Render(w, r, httperr.ErrInvalidRequest(err))
			return
		}

		// Fetch people info from external API
		people, err := fetchPeopleInfo(passport.Serie, passport.Number, extAPIPort)
		if err != nil {
			log.Error("failed to fetch people info", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		// Now you have `people` containing the data retrieved from the external API
		log.Debug("fetched people info", slog.Any("people", people))

		user := models.User{
			People:   people,
			Passport: passport,
		}
		userId, err := userCreator.CreateUser(r.Context(), user)
		if err != nil {
			log.Error("failed to create user", l.Err(err))

			render.Render(w, r, httperr.ErrInternal(err))
			return
		}

		log.Info("created user", slog.Int("user_id", int(userId)))

		resp := CreateUserResponse{UserID: userId}
		render.Status(r, http.StatusCreated)
		render.JSON(w, r, resp)
	}
}

// Making call to external API to fetch data for user who matches the given passport data
func fetchPeopleInfo(passportSerie, passportNumber string, extAPIPort int) (models.People, error) {
	url := fmt.Sprintf("http://localhost:%d/info?passportSerie=%s&passportNumber=%s", extAPIPort, passportSerie, passportNumber)

	resp, err := http.Get(url)
	if err != nil {
		return models.People{}, fmt.Errorf("failed to fetch people info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.People{}, fmt.Errorf("unexpected status code: %v", resp.Status)
	}

	var people models.People
	if err := json.NewDecoder(resp.Body).Decode(&people); err != nil {
		return models.People{}, fmt.Errorf("failed to decode response body: %v", err)
	}

	return people, nil
}
