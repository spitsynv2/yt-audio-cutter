package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
	Id         string    `json:"id"`
	Name       string    `json:"name" binding:"required"`
	YoutubeURL string    `json:"youtube_url" binding:"required,url"`
	StartTime  *Duration `json:"start_time" binding:"required"` // e.g. "10s"
	EndTime    *Duration `json:"end_time" binding:"required"`
	Status     JobStatus `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	FileUrl    string    `json:"file_url"`
}

type Duration struct {
	time.Duration
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Duration.String())
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	dur, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	d.Duration = dur
	return nil
}

func (d *Duration) Scan(value interface{}) error {
	if value == nil {
		return nil
	}

	var s string
	switch v := value.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	case int64:
		d.Duration = time.Duration(v) * time.Second
		return nil
	default:
		return fmt.Errorf("unsupported type %T for Duration", v)
	}

	if dur, err := time.ParseDuration(s); err == nil {
		d.Duration = dur
		return nil
	}

	parts := strings.Split(s, ":")
	if len(parts) == 3 {
		h, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		sec, _ := strconv.ParseFloat(parts[2], 64)

		total := time.Duration(h)*time.Hour +
			time.Duration(m)*time.Minute +
			time.Duration(sec*float64(time.Second))

		d.Duration = total
		return nil
	}

	return errors.New("invalid duration format: " + s)
}

func (d Duration) Value() (driver.Value, error) {
	return d.Duration.String(), nil
}

func FormatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
