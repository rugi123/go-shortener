package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rugi123/go-shortener/internal/config"
)

type Url struct {
	Id    int
	Url   string
	Alias string
}

func InintDB(ctx context.Context, cfg config.PostgresConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к db: %s", err)
	}
	return conn, nil
}

func SearchUrlFromDB(alias string, ctx context.Context, conn pgx.Conn) (*Url, error) {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфига: %s", err)
	}

	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT id, url, alias FROM %s WHERE alias = $1", cfg.Postgres.TableName), alias)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска в db: %s", err)
	}

	defer rows.Close()
	for rows.Next() {
		url := Url{}
		err := rows.Scan(&url.Id, &url.Url, &url.Alias)
		if err != nil {
			return nil, fmt.Errorf("ошибка скана db: %s", err)
		}
		if url.Alias == alias {
			return &url, nil
		}
	}
	return nil, err
}
func AddUrlToDB(base_url string, alias string, ctx context.Context, conn pgx.Conn) error {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфига: %s", err)
	}

	url, err := SearchUrlFromDB(alias, ctx, conn)
	if err != nil {
		return fmt.Errorf("ошибка поиска в db: %s", err)
	}
	if url != nil {
		return fmt.Errorf("найден такой же alias")
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("INSERT INTO %s (url, alias) VALUES ($1, $2)", cfg.Postgres.TableName), base_url, alias)
	if err != nil {
		return fmt.Errorf("ошибка вставки в db: %s", err)
	}
	return err
}
func RemoveUrlFromDB() error {
	return nil
}
