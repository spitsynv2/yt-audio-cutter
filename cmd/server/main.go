package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spitsynv2/yt-audio-cutter/internal/api"
)

func main() {
	router := gin.Default()
	api.RegisterRoutes(router)

	log.Fatal(router.Run(":8080")) // Listen and serve
}
