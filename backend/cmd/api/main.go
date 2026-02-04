package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"pack-calculator/internal/adapters/cache"
	"pack-calculator/internal/adapters/repository"
	"pack-calculator/internal/app"
	"pack-calculator/internal/config"
	httptransport "pack-calculator/internal/transport/http"
	"pack-calculator/pkg/logger"
)

func main() {
	log := logger.Default()

	cfg, err := config.Load()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	repo, err := repository.NewPostgresRepository(cfg.DB.DSN())
	if err != nil {
		log.Error("Failed to initialize repository", "error", err)
		os.Exit(1)
	}
	defer repo.Close()

	redisCache, err := cache.NewRedisCache(cfg.Redis.Addr(), cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		log.Error("Failed to initialize cache", "error", err)
		os.Exit(1)
	}
	defer redisCache.Close()

	packService := app.NewPackService(repo, redisCache)
	handler := httptransport.NewHandler(packService)
	router := httptransport.SetupRoutes(handler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	go func() {
		log.Info("Server starting", "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Server failed", "error", err)
			os.Exit(1)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
		os.Exit(1)
	}

	log.Info("Server exited")
}
