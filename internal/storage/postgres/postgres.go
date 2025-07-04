package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rugi123/go-shortener/internal/config"
)

func InintDB(ctx context.Context, cfg config.PostgresConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к db:", err)
	}
	return conn, nil
}
