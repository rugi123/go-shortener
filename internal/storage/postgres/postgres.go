package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/rugi123/go-shortener/internal/config"
	"github.com/rugi123/go-shortener/internal/domain/model"
)

func InintDB(ctx context.Context, cfg config.PostgresConfig) (*pgx.Conn, error) {
	conn, err := pgx.Connect(ctx, cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к db: %s", err)
	}
	return conn, nil
}

func SearchLink(short_url string, ctx context.Context, conn pgx.Conn) (*model.Link, error) {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки конфига: %s", err)
	}

	rows, err := conn.Query(ctx, fmt.Sprintf("SELECT id, original_url, short_url FROM %s WHERE short_url = $1", cfg.Postgres.TableName), short_url)
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска в db: %s", err)
	}

	defer rows.Close()
	for rows.Next() {
		url := model.Link{}
		err := rows.Scan(&url.ID, &url.OriginalUrl, &url.ShortUrl)
		if err != nil {
			return nil, fmt.Errorf("ошибка скана db: %s", err)
		}
		if url.ShortUrl == short_url {
			return &url, nil
		}
	}
	return nil, err
}
func SaveLink(link model.Link, ctx context.Context, conn pgx.Conn) error {
	cfg, err := config.Load("internal/config/config.yaml")
	if err != nil {
		return fmt.Errorf("ошибка загрузки конфига: %s", err)
	}

	url, err := SearchLink(link.ShortUrl, ctx, conn)
	if err != nil {
		return fmt.Errorf("ошибка поиска в db: %s", err)
	}
	if url != nil {
		return fmt.Errorf("найден такой же alias")
	}
	_, err = conn.Exec(ctx, fmt.Sprintf("INSERT INTO %s (original_url, short_url) VALUES ($1, $2)",
		cfg.Postgres.TableName),
		link.OriginalUrl,
		link.ShortUrl)
	if err != nil {
		return fmt.Errorf("ошибка вставки в db: %s", err)
	}
	return err
}
func DeleteLink() error {
	return nil
}
