package configs

import (
	"time"
	postgres "workwithimages/pkg/pgx"
)

var Psql postgres.Postgres = postgres.Postgres{
	Host:     "localhost",
	Port:     "5432",
	Username: "postgres",
	Password: "1",
}

type HTTP struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var HOST = HTTP{
	Host:         "127.0.0.1",
	Port:         ":8080",
	ReadTimeout:  10 * time.Second,
	WriteTimeout: 10 * time.Second,
}
