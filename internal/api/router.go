package api

import "github.com/gin-gonic/gin"

func RegisterRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"server": "running"})
	})

	router.POST("/jobs", CreateJob)
	router.GET("/jobs/:id", GetJobById)
	router.GET("/jobs", GetJobs)
}
