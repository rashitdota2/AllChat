package postgres

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

type Postgres struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	DBName   string `mapstructure:"dbname"`
}

func NewPostgresDB(cfg Postgres) (*pgxpool.Pool, error) {
	dataSourceName := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := pgxpool.Connect(ctx, dataSourceName)
	if err != nil {
		return nil, err
	}

	// her gezek taze context bolsa gowy
	if err = db.Ping(context.Background()); err != nil {
		return nil, errors.New("db.Ping")
	}

	return db, nil
}
