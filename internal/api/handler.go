package api

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spitsynv2/yt-audio-cutter/internal/model"
	"net/http"
	"time"
)

var inMemoryJobs = make(map[string]model.Job) // temporary in-memory storage

func CreateJob(c *gin.Context) {
	var jobInput model.Job

	if err := c.ShouldBindJSON(&jobInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate job
	jobID := uuid.NewString()
	newJob := model.Job{
		ID:         jobID,
		YoutubeURL: jobInput.YoutubeURL,
		StartTime:  jobInput.StartTime,
		EndTime:    jobInput.EndTime,
		Status:     model.StatusPending,
		CreatedAt:  time.Now(),
	}

	inMemoryJobs[jobID] = newJob

	c.JSON(http.StatusCreated, gin.H{"id": jobID})
}

func GetJobById(c *gin.Context) {
	id := c.Param("id")

	job, exists := inMemoryJobs[id]

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context) {
	var jobs = make([]model.Job, 0, len(inMemoryJobs))

	for _, job := range inMemoryJobs {
		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, jobs)
}
