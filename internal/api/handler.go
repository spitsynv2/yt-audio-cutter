package api

import (
	"encoding/json"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spitsynv2/yt-audio-cutter/internal/model"
	"github.com/spitsynv2/yt-audio-cutter/internal/service"
	"github.com/spitsynv2/yt-audio-cutter/internal/store"
)

var inMemoryStore = &store.MemoryJobStore{Jobs: make(map[string]model.Job)}

func GetStartPage(c *gin.Context) {
	endpoints := []struct {
		Name     string `json:"name"`
		Endpoint string `json:"endpoint"`
	}{
		{"GET server status", "/health"},
		{"GET jobs", "/jobs"},
		{"GET job by id", "/jobs/id"},
		{"POST create job", "/jobs"},
	}

	prettyJSON, err := json.MarshalIndent(endpoints, "", "  ")
	if err != nil {
		c.String(http.StatusInternalServerError, "JSON marshal error: %s", err)
		return
	}

	c.Data(http.StatusOK, "application/json", prettyJSON)
}

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

	inMemoryStore.PutJob(newJob)
	service.ProcessJob(jobID, inMemoryStore)

	c.JSON(http.StatusCreated, gin.H{"id": jobID})
}

func GetJobById(c *gin.Context) {
	id := c.Param("id")

	job, exists := inMemoryStore.GetJob(id)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context) {
	var jobs = make([]model.Job, 0, len(inMemoryStore.Jobs))

	for _, job := range inMemoryStore.Jobs {
		jobs = append(jobs, job)
	}

	c.JSON(http.StatusOK, jobs)
}

func UpdateJob(c *gin.Context) {
	id := c.Param("id")

	job, exists := inMemoryStore.GetJob(id)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	if job.Status != model.StatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot apdate not pending jobs"})
		return
	}

	var updateInput struct {
		StartTime *model.Duration `json:"start_time" binding:"required"`
		EndTime   *model.Duration `json:"end_time" binding:"required"`
	}

	if err := c.ShouldBindJSON(&updateInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job.StartTime = updateInput.StartTime
	job.EndTime = updateInput.EndTime
	inMemoryStore.Update(id, job)
	c.JSON(http.StatusOK, job)
}

func DeleteJob(c *gin.Context) {
	id := c.Param("id")

	job, exists := inMemoryStore.GetJob(id)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}

	if job.Status != model.StatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot delete job once processing has started"})
		return
	}

	inMemoryStore.DeleteJob(id)
	c.Status(http.StatusNoContent)
}
