package pg_client

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ConfigPG struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
	TimeZone string
}

func NewPostgresDB(cfg ConfigPG) (*pgxpool.Pool, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)
	ctx := context.Background()
	db, err := pgxpool.Connect(ctx, dsn)

	if err != nil {
		return nil, err
	}
	return db, nil
}
