package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/senyabanana/avito-shop-service/internal/config"
)

func NewPostgresDB(ctx context.Context, cfg *config.Config) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.PostgresHost, cfg.PostgresPort, cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB, cfg.SSLMode)

	// Таймаут на все подключение (включая открытие соединения и пинг)
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	db, err := sqlx.ConnectContext(connectCtx, "postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return db, nil
}
