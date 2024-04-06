package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"log"
)

type service struct {
	Repo
}

func NewService(repository Repo) Service {
	return &service{
		repository,
	}
}

func (s *service) IncVotes(ctx context.Context, id int64, chat *Chat) {
	ms, err := s.Repo.IncVotesByID(ctx, id)
	log.Printf("IncVotesByID: %v", ms)
	if err != nil {
		return
	}

	chat.Broadcast <- ms
}

func (s *service) DecVotes(ctx context.Context, id int64, chat *Chat) {
	ms, err := s.Repo.DecVotesByID(ctx, id)
	log.Printf("DecVotesByID: %v", ms)
	if err != nil {
		return
	}

	chat.Broadcast <- ms
}

func (s *service) WriteMessage(ctx context.Context, cl *Client) {
	defer func() {
		cl.Conn.Close()
	}()

	for {
		_, ok := <-cl.Message
		if !ok {
			return
		}

		ms, err := s.Repo.GetMostRecentMsg(ctx)
		if err != nil {
			return
		}

		err = cl.Conn.WriteJSON(ms)
		if err != nil {
			return
		}
	}
}

func (s *service) ReadMessage(ctx context.Context, chat *Chat, cl *Client) {
	defer func() {
		chat.Unregister <- cl
		cl.Conn.Close()
	}()

	for {
		_, m, err := cl.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message := &Message{
			MessageID: 0,
			Content:   string(m),
			ChatID:    cl.ChatID,
			Username:  cl.Username,
			Votes:     0,
		}

		ms, err := s.Repo.CreateMessage(ctx, message)
		if err != nil {
			return
		}

		chat.Broadcast <- ms

	}
}
