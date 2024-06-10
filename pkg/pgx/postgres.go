package pgx

import "database/sql"

func NewDB() *sql.DB {
	db, _ := sql.Open("pgx", "psql")
	return db
}
