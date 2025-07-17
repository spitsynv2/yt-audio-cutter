package store

import "github.com/spitsynv2/yt-audio-cutter/internal/model"

type JobStorage interface {
	PutJob(job model.Job)
	GetJob(id string) (model.Job, bool)
	Update(id string, new model.Job) error
	DeleteJob(id string) error
}
