package main

import (
	"context"
	"fmt"

	"github.com/rugi123/go-shortener/internal/config"
	"github.com/rugi123/go-shortener/internal/storage/postgres"
)

func main() {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		fmt.Println("ошибка загрузки конфига:", err)
	}
	conn, err := postgres.InintDB(context.Background(), cfg.Postgres)
	if err != nil {
		fmt.Println("ошибка инициализации db:", err)
	}
	fmt.Println(conn)
}
