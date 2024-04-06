package user

import (
	"context"
	"database/sql"
)

type DBOps interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type repo struct {
	db DBOps
}

func NewRepo(db DBOps) Repo {
	return &repo{db: db}
}

func (r *repo) CreateUser(ctx context.Context, user *User) (*User, error) {
	var lastInsertId int
	query := "INSERT INTO users(username, password, first_name, last_name) VALUES ($1, $2, $3, $4) returning id"
	err := r.db.QueryRowContext(ctx, query, user.Username, user.Password, user.Firstname, user.Lastname).Scan(&lastInsertId)
	if err != nil {
		return &User{}, err
	}

	user.ID = int64(lastInsertId)
	return user, nil
}

func (r *repo) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	u := User{}

	query := "SELECT id, username, password, first_name, last_name FROM users WHERE username = $1"
	err := r.db.QueryRowContext(ctx, query, username).Scan(&u.ID, &u.Username, &u.Password, &u.Firstname, &u.Lastname)
	if err != nil {
		return &User{}, nil
	}

	return &u, nil
}
