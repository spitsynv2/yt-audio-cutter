package main

import (
	"context"
	"log"
	"net/http"
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
	router.Use(api.TimeoutMiddleware(3 * time.Minute))
	api.RegisterRoutes(router)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Println("server started on :8080")

	<-root.Done()
	log.Println("shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server exiting")
}
