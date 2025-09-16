package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spitsynv2/yt-audio-cutter/internal/api"
	"github.com/spitsynv2/yt-audio-cutter/internal/store"
)

func main() {
	root, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err := store.InitConnection(root); err != nil {
		log.Fatal(err)
	}
	defer store.DB.Close()

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(api.TimeoutMiddleware(15 * time.Second))

	api.RegisterRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
