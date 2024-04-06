package websocket

import (
	"context"
	"database/sql"
	"log"
)

type DBOps interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

type Repo interface {
	CreateMessage(ctx context.Context, msg *Message) (*Message, error)
	GetMsgById(ctx context.Context, id int) (*Message, error)
	GetMostRecentMsg(ctx context.Context) (*Message, error)
	IncVotesByID(ctx context.Context, id int64) (*Message, error)
	DecVotesByID(ctx context.Context, id int64) (*Message, error)
}

type repo struct {
	db DBOps
}

func (r *repo) CreateMessage(ctx context.Context, msg *Message) (*Message, error) {
	var lastInsertId int
	query := "INSERT INTO messages(content, chat_id, username, votes) VALUES ($1, $2, $3, $4) returning message_id"
	err := r.db.QueryRowContext(ctx, query, msg.Content, msg.ChatID, msg.Username, msg.Votes).Scan(&lastInsertId)
	if err != nil {
		log.Fatalf("CreateMessage Query failed: %s", err)
	}

	msg.MessageID = int64(lastInsertId)
	log.Printf("Create message with id: %d", msg.MessageID) //fixme: remove later
	return msg, nil
}

func (r *repo) GetMsgById(ctx context.Context, id int) (*Message, error) {
	m := Message{}

	query := "SELECT message_id, username, chat_id, content, votes FROM messages WHERE message_id = $1"
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.MessageID, &m.Username, &m.ChatID, &m.Content, &m.Votes)
	if err != nil {
		log.Fatalf("GetMsgById Query failed: %s", err)
	}

	return &m, nil
}

func (r *repo) GetMostRecentMsg(ctx context.Context) (*Message, error) {
	m := Message{}

	query := "SELECT message_id, username, chat_id, content, votes FROM messages ORDER BY created_at DESC LIMIT 1"
	err := r.db.QueryRowContext(ctx, query).Scan(&m.MessageID, &m.Username, &m.ChatID, &m.Content, &m.Votes)
	if err != nil {
		log.Fatalf("GetMostRecentMsg Query failed: %s", err)
	}

	return &m, nil
}

func (r *repo) IncVotesByID(ctx context.Context, id int64) (*Message, error) {
	m := Message{}

	query := `UPDATE messages SET votes = votes + 1, created_at = now() WHERE message_id = $1 RETURNING message_id, username, chat_id, content, votes`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.MessageID, &m.Username, &m.ChatID, &m.Content, &m.Votes)
	if err != nil {
		log.Fatalf("IncVotesByID Query failed: %s", err)
	}

	return &m, nil
}

func (r *repo) DecVotesByID(ctx context.Context, id int64) (*Message, error) {
	m := Message{}

	query := `UPDATE messages SET votes = votes - 1, created_at = now() WHERE message_id = $1 RETURNING message_id, username, chat_id, content, votes`
	err := r.db.QueryRowContext(ctx, query, id).Scan(&m.MessageID, &m.Username, &m.ChatID, &m.Content, &m.Votes)
	if err != nil {
		log.Fatalf("DecVotesByID Query failed: %s", err)
	}

	return &m, nil
}

func NewRepo(db DBOps) Repo {
	return &repo{db: db}
}
