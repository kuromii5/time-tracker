package repo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/kuromii5/time-tracker/internal/models"
	l "github.com/kuromii5/time-tracker/pkg/logger"
)

var (
	ErrAlreadyDone = errors.New("worklog was already finished")
)

func (db *DB) StartWorklog(ctx context.Context, task string, userID int32) (int32, error) {
	query := `
		INSERT INTO worklogs (user_id, task, started_at)
		VALUES ($1, $2, NOW())
		RETURNING ID
	`
	log := db.log.With(slog.String("task", task), slog.Int("user_id", int(userID)))
	log.Debug("executing query", slog.String("query", query))

	var worklogId int32
	err := db.pool.QueryRow(ctx, query, userID, task).Scan(&worklogId)
	if err != nil {
		log.Error("failed to execute query", l.Err(err))

		return 0, fmt.Errorf("%s: %w", "repo.StartWorklog", err)
	}

	log.Debug("worklog started successfully", slog.Int("worklog_id", int(worklogId)))

	return worklogId, nil
}

func (db *DB) FinishWorklog(ctx context.Context, worklogID int32) error {
	query := `
		UPDATE worklogs
		SET finished_at = NOW()
		WHERE id = $1 AND finished_at IS NULL
		RETURNING id
	`
	log := db.log.With(slog.Int("worklog_id", int(worklogID)))
	log.Debug("executing query", slog.String("query", query))

	var id int32
	err := db.pool.QueryRow(ctx, query, worklogID).Scan(&id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// No rows were updated, meaning the worklog was already finished
			return ErrAlreadyDone
		}
		log.Error("failed to execute query", l.Err(err))

		return fmt.Errorf("%s: %w", "repo.FinishWorklog", err)
	}

	log.Debug("worklog finished successfully")

	return nil
}

func (db *DB) Worklogs(ctx context.Context, userID int32, startDate, endDate time.Time) ([]models.Worklog, error) {
	query := `
		SELECT * FROM worklogs
		WHERE user_id = $1 AND started_at >= $2 AND (finished_at <= $3 OR finished_at IS NULL)
		ORDER BY duration DESC
	`
	log := db.log.With(slog.Int("user_id", int(userID)), slog.Time("start_date", startDate), slog.Time("end_date", endDate))
	log.Debug("executing query", slog.String("query", query))

	rows, err := db.pool.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		log.Error("failed to execute query", l.Err(err))

		return nil, fmt.Errorf("%s: %w", "repo.Worklogs", err)
	}
	defer rows.Close()

	var worklogs []models.Worklog
	for rows.Next() {
		var worklog models.Worklog
		err := rows.Scan(&worklog.ID, &worklog.UserID, &worklog.StartedAt, &worklog.FinishedAt, &worklog.Task, &worklog.Duration)
		if err != nil {
			log.Error("failed to scan row", l.Err(err))

			return nil, fmt.Errorf("%s: %w", "repo.Worklogs", err)
		}

		worklogs = append(worklogs, worklog)
	}
	if err := rows.Err(); err != nil {
		log.Error("rows error", l.Err(err))

		return nil, fmt.Errorf("%s: %w", "repo.Worklogs", err)
	}

	log.Debug("worklogs retrieved successfully", slog.Int("count", len(worklogs)))

	return worklogs, nil
}
