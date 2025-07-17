package store

import "github.com/spitsynv2/yt-audio-cutter/internal/model"

type MemoryJobStore struct {
	Jobs map[string]model.Job
}

func (m *MemoryJobStore) PutJob(job model.Job) {
	m.Jobs[job.ID] = job
}

func (m *MemoryJobStore) GetJob(id string) (model.Job, bool) {
	job, ok := m.Jobs[id]
	return job, ok
}

func (m *MemoryJobStore) Update(id string, updated model.Job) error {
	_, ok := m.Jobs[id]
	if !ok {
		return nil
	}

	m.Jobs[id] = updated
	return nil
}

func (m *MemoryJobStore) DeleteJob(id string) error {
	delete(m.Jobs, id)
	return nil
}
