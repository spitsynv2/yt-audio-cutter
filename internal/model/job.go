package model

import (
	"strings"
	"time"
)

type JobStatus string

const (
	StatusPending JobStatus = "pending"
	StatusRunning JobStatus = "running"
	StatusDone    JobStatus = "done"
	StatusFailed  JobStatus = "failed"
)

type Job struct {
	ID         string    `json:"id"`
	YoutubeURL string    `json:"youtube_url" binding:"required,url"`
	StartTime  *Duration `json:"start_time" binding:"required"` // e.g. "10s"
	EndTime    *Duration `json:"end_time" binding:"required"`
	Status     JobStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}

type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = duration
	return nil
}
