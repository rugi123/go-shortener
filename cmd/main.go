package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rugi123/go-shortener/internal/config"
	"github.com/rugi123/go-shortener/internal/handlers"
	"github.com/rugi123/go-shortener/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		fmt.Println("ошибка загрузки конфига: ", err)
	}
	conn, err := postgres.InintDB(context.Background(), cfg.Postgres)
	if err != nil {
		fmt.Println("ошибка инициализации db: ", err)
	}

	r := gin.Default()

	r.POST("/shorten", func(c *gin.Context) {
		var data map[string]string

		if err := c.ShouldBindJSON(&data); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		fmt.Println(data)
		message, exists := data["url"]
		if !exists {
			c.JSON(http.StatusBadRequest, gin.H{"error": "field 'url' is required"})
			return
		}

		short_url, err := handlers.Shortener(message, r, conn, c)
		if err != nil {
			fmt.Println("ошибка создания алиаса: ", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"short_url": short_url,
		})
	})
	r.GET("/long_url", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"tralalelo": "tralala",
		})
	})

	r.Run(":" + cfg.App.Port)
}
