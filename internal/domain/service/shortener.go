package service

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/rugi123/go-shortener/internal/config"
)

func GenerateShortURL() (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		return "", fmt.Errorf("ошибка загрузки конфига: %s", err)
	}

	b := make([]byte, cfg.App.URLLength)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}

	return "/" + string(b), nil
}
