package api

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.StaticFile("/", "./ui/index.html")

	router.POST("/jobs", CreateJob)
	router.GET("/jobs/:id", GetJobById)
	router.GET("/jobs/process/:id", ProcessJobId)
	router.GET("/jobs", GetJobs)
	router.PUT("/jobs/:id", UpdateJob)
	router.DELETE("/jobs/:id", DeleteJob)
}
