package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohan7-code/url-shortener/config"
	"github.com/mohan7-code/url-shortener/database"
	"github.com/mohan7-code/url-shortener/routes"
	"github.com/mohan7-code/url-shortener/utils/cache"
)

func main() {
	log.Printf("URL Shortener Service starting...")

	cnf, err := config.LoadConfig(".env")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	database.Init(&database.Config{
		URL:       cnf.DatabaseUrl,
		MaxDBConn: cnf.MaxDBConn,
	})
	cache.SetRedis()
	r := routes.GetRouter()

	server := &http.Server{
		Addr:    ":" + cnf.ServerPort,
		Handler: r,
	}
	go func() {
		log.Printf("Server running on port %s", cnf.ServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Shutdown failed: %v", err)
	}

	sqlDB, _ := database.DB.DB()
	sqlDB.Close()

	log.Println("Server stopped gracefully.")
}
