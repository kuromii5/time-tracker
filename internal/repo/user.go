package repo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5"
	"github.com/kuromii5/time-tracker/internal/models"
	"github.com/kuromii5/time-tracker/internal/utils"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrPassportDuplicate = errors.New("user with such serie and number already exists")
)

func (db *DB) CreateUser(ctx context.Context, user models.User) (int32, error) {
	query := `
		INSERT INTO users (passport_serie, passport_number, name, surname, patronymic, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id
	`
	log := db.log.With(slog.Any("user", user))
	log.Debug("executing query", slog.String("query", query))

	var userId int32
	err := db.pool.QueryRow(ctx, query,
		user.Passport.Serie, user.Passport.Number, user.People.Name, user.People.Surname, user.People.Patronymic, user.People.Address).
		Scan(&userId)
	if err != nil {
		// Check if the error is a unique constraint violation
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is the unique_violation error code in PostgreSQL
			log.Error("user with such serie and number already exists", l.Err(ErrPassportDuplicate))

			return 0, ErrPassportDuplicate
		}
		log.Error("failed to execute query", l.Err(err))

		return 0, fmt.Errorf("%s: %w", "repo.CreateUser", err)
	}

	log.Debug("successfully created user", slog.Int("user_id", int(userId)))

	return userId, nil
}

func (db *DB) Users(ctx context.Context, filter models.FilterBy, settings models.Pagination) ([]models.User, error) {
	query, args := utils.BuildGetUsersQuery(filter, settings)

	log := db.log.With(slog.Any("filter", filter), slog.Any("pagination", settings))
	log.Debug("executing query", slog.String("query", query), slog.Any("args", args))

	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute query", l.Err(err))

		return nil, fmt.Errorf("%s: %w", "repo.Users", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt, &user.Passport.Serie, &user.Passport.Number, &user.People.Name, &user.People.Surname, &user.People.Patronymic, &user.People.Address)
		if err != nil {
			log.Error("failed to scan row", l.Err(err))

			return nil, fmt.Errorf("%s: %w", "repo.Users", err)
		}

		users = append(users, user)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows error", l.Err(err))

		return nil, fmt.Errorf("%s: %w", "repo.Users", err)
	}

	log.Debug("successfully retrieved users")

	return users, nil
}

func (db *DB) DeleteUser(ctx context.Context, id int32) error {
	query := "DELETE FROM users WHERE id = $1"

	log := db.log.With(slog.Int("user_id", int(id)))
	log.Debug("executing query", slog.String("query", query))

	_, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Error("user not found", l.Err(err))

			return fmt.Errorf("%s: %w", "repo.DeleteUser", ErrUserNotFound)
		}
		log.Error("failed to execute query", l.Err(err))

		return fmt.Errorf("%s: %w", "repo.DeleteUser", err)
	}

	log.Debug("successfully deleted user")

	return nil
}

func (db *DB) UpdateUser(ctx context.Context, user models.User) error {
	query, args := utils.BuildUpdateUserQuery(user)

	log := db.log.With(slog.Any("user", user))
	log.Debug("executing query", slog.String("query", query), slog.Any("args", args))

	_, err := db.pool.Exec(ctx, query, args...)
	if err != nil {
		log.Error("failed to execute query", l.Err(err))

		return fmt.Errorf("%s: %w", "repo.UpdateUser", err)
	}

	log.Debug("successfully updated user")

	return nil
}
