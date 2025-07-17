package service

import (
	"log"
	"time"

	"github.com/spitsynv2/yt-audio-cutter/internal/model"
	"github.com/spitsynv2/yt-audio-cutter/internal/store"
)

// Simulate processing in background
func ProcessJob(id string, store store.JobStorage) {
	go func() {
		log.Printf("Job %s: queued for processing...", id)

		time.Sleep(2 * time.Second) // simulate startup delay

		job, exists := store.GetJob(id)
		if !exists {
			log.Printf("Job %s: not found in storage", id)
			return
		}

		job.Status = model.StatusRunning
		store.Update(id, job)
		log.Printf("Job %s: running", id)

		time.Sleep(5 * time.Second) // simulate processing

		job.Status = model.StatusDone
		store.Update(id, job)
		log.Printf("Job %s: done", id)
	}()
}
