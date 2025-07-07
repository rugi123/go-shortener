package main

import (
	"context"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/rugi123/go-shortener/internal/config"
	"github.com/rugi123/go-shortener/internal/domain/service"
	"github.com/rugi123/go-shortener/internal/handlers"
	"github.com/rugi123/go-shortener/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		log.Fatal("failed to load config: ", err)
	}
	storage, err := postgres.NewPGStorage(context.Background(), &cfg.Postgres)
	if err != nil {
		log.Fatal("failed to create new pg storage: ", err)
	}
	defer storage.Close()

	service := service.NewShortenerService(storage, cfg.App.URLLength)

	handler := handlers.NewShortenHandler(service, "http://localhost:"+cfg.App.Port)

	router := gin.Default()

	router.POST("/api/shorten", handler.Shorten)
	router.GET("/:key", handler.Redirect)

	router.Run()
}
