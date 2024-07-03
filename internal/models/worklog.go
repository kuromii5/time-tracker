package models

import "time"

type Worklog struct {
	ID         int32         `json:"id"`
	UserID     int32         `json:"user_id"`
	Task       string        `json:"task"`
	StartedAt  time.Time     `json:"start_time"`
	FinishedAt time.Time     `json:"end_time"`
	Duration   time.Duration `json:"duration"`
}
