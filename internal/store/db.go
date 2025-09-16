package store

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/spitsynv2/yt-audio-cutter/config"
	"github.com/spitsynv2/yt-audio-cutter/internal/model"
)

const (
	createJobsQuery = "INSERT INTO yt_convertor.jobs(id, youtube_url, start_time, end_time, status) VALUES($1, $2, $3, $4, $5)"
	getJobQuery     = "SELECT id, youtube_url, start_time, end_time, status, created_at FROM yt_convertor.jobs WHERE id=$1"
	getAllJobsQuery = "SELECT id, youtube_url, start_time, end_time, status, created_at FROM yt_convertor.jobs"
	updateJobQuery  = "UPDATE yt_convertor.jobs SET youtube_url = $1, start_time = $2, end_time = $3, status = $4 WHERE id = $5"
	deleteJobQuery  = "DELETE FROM yt_convertor.jobs WHERE id = $1"
)

var DB *sql.DB

type JobStorage interface {
	CreateJob(ctx context.Context, job model.Job)
	GetJob(ctx context.Context, id string) (model.Job, bool)
	UpdateJob(ctx context.Context, id string, new model.Job) error
	DeleteJob(ctx context.Context, id string) error
}

func InitConnection(ctx context.Context) error {
	db, err := sql.Open("pgx", config.Conf.Dsn)
	if err != nil {
		return err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(30 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return err
	}
	DB = db

	return nil
}

func CreateJob(ctx context.Context, job model.Job) (int64, error) {
	res, err := DB.ExecContext(ctx, createJobsQuery, job.Id, job.YoutubeURL, job.StartTime, job.EndTime, job.Status)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func GetJob(ctx context.Context, id string) (model.Job, error) {
	var job model.Job
	err := DB.QueryRowContext(ctx, getJobQuery, id).Scan(&job.Id, &job.YoutubeURL, &job.StartTime, &job.EndTime, &job.Status, &job.CreatedAt)
	return job, err
}

func GetJobs(ctx context.Context) ([]model.Job, error) {
	rows, err := DB.QueryContext(ctx, getAllJobsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var jobs []model.Job
	for rows.Next() {
		var job model.Job
		if err := rows.Scan(&job.Id, &job.YoutubeURL, &job.StartTime, &job.EndTime, &job.Status, &job.CreatedAt); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func UpdateJob(ctx context.Context, job model.Job) (int64, error) {
	res, err := DB.ExecContext(ctx, updateJobQuery, job.YoutubeURL, job.StartTime, job.EndTime, job.Status, job.Id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func DeleteJob(ctx context.Context, id string) (int64, error) {
	res, err := DB.ExecContext(ctx, deleteJobQuery, id)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}
