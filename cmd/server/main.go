package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"asset-core/internal/api"
	"asset-core/internal/config"
	"asset-core/internal/infrastructure/kafka"
	"asset-core/internal/infrastructure/mysql"
	"asset-core/internal/infrastructure/redis"
	"asset-core/internal/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.App.Env)
	defer log.Sync()

	db := mysql.New(cfg.MySQL)
	rdb := redis.New(cfg.Redis)
	producer := kafka.NewProducer(cfg.Kafka, log)
	defer producer.Close()

	router := api.NewRouter(api.Dependencies{
		Config:        cfg,
		Logger:        log,
		DB:            db,
		Redis:         rdb,
		EventProducer: producer,
	})

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", cfg.App.Port),
		Handler:           router,
		ReadHeaderTimeout: 10 * time.Second,
	}

	go func() {
		log.Info("server starting", logger.Int("port", cfg.App.Port))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("server failed", logger.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server shutdown failed", logger.Error(err))
	}
	log.Info("server stopped")
}
