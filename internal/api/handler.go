package api

import (
	"database/sql"
	"encoding/json"

	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spitsynv2/yt-audio-cutter/internal/model"
	"github.com/spitsynv2/yt-audio-cutter/internal/store"
)

func TimeoutMiddleware(d time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), d)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusGatewayTimeout, gin.H{"error": "request timed out"})
			c.Abort()
		}
	}
}

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

	jobID := uuid.NewString()
	newJob := model.Job{
		Id:         jobID,
		YoutubeURL: jobInput.YoutubeURL,
		StartTime:  jobInput.StartTime,
		EndTime:    jobInput.EndTime,
		Status:     model.StatusPending,
		CreatedAt:  time.Now(),
	}

	_, err := store.CreateJob(c.Request.Context(), newJob)
	if err != nil {
		c.String(http.StatusInternalServerError, "Job creation error: %s", err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": jobID})
}

func GetJobById(c *gin.Context) {
	id := c.Param("id")

	job, err := store.GetJob(c.Request.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func GetJobs(c *gin.Context) {

	jobs, err := store.GetJobs(c.Request.Context())
	if errors.Is(err, sql.ErrNoRows) {
		c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, jobs)
}

func UpdateJob(c *gin.Context) {
	id := c.Param("id")

	job, err := store.GetJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}

	if job.Status != model.StatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot apdate active jobs"})
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

	_, err = store.UpdateJob(c.Request.Context(), job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, job)
}

func DeleteJob(c *gin.Context) {
	id := c.Param("id")

	job, err := store.GetJob(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"No such job error": err})
		return
	}

	if job.Status != model.StatusPending {
		c.JSON(http.StatusConflict, gin.H{"error": "cannot delete job once processing has started"})
		return
	}

	_, err = store.DeleteJob(c.Request.Context(), job.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusOK)
}
