package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rugi123/go-shortener/internal/config"
	"github.com/rugi123/go-shortener/internal/domain/model"
)

type PGStorage struct {
	pool      *pgxpool.Pool
	tableName string
}

func NewPGStorage(ctx context.Context, cfg *config.PostgresConfig) (*PGStorage, error) {
	poolConfig, err := pgxpool.ParseConfig(cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.HealthCheckPeriod = 1 * time.Minute

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &PGStorage{
		pool:      pool,
		tableName: cfg.TableName,
	}, nil
}

func (s *PGStorage) SaveLink(ctx context.Context, link *model.Link) error {
	query := fmt.Sprintf(`
		INSERT INTO %s (original_url, short_key) 
		VALUES ($1, $2)
		ON CONFLICT (short_key) DO NOTHING`,
		s.tableName)
	_, err := s.pool.Exec(ctx, query, link.OriginalURL, link.ShortKey)
	if err != nil {
		return fmt.Errorf("failed to exec query: %w", err)
	}

	return nil
}

func (s *PGStorage) GetLinkByKey(ctx context.Context, key string) (*model.Link, error) {
	query := fmt.Sprintf(`
		SELECT id, original_url, short_key 
		FROM %s 
		WHERE short_key = $1`,
		s.tableName)

	var link model.Link
	err := s.pool.QueryRow(ctx, query, key).Scan(&link.ID, &link.OriginalURL, &link.ShortKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get link: %w", err)
	}

	return &link, nil
}

func (s *PGStorage) Close() {
	s.pool.Close()
}
