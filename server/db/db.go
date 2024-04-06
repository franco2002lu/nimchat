package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDB() (*Database, error) {
	db, err := sql.Open("postgres", "postgresql://root:franco0808@localhost:5433/nimble-chat?sslmode=disable")
	if err != nil {
		return nil, err
	}

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) GetDB() *sql.DB {
	return d.db
}
