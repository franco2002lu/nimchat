package websocket

import (
	"context"
	"github.com/gorilla/websocket"
	"sync"
)

type Message struct {
	MessageID int64  `json:"message_id"`
	Content   string `json:"content"`
	ChatID    int64  `json:"chat_id"`
	Username  string `json:"username"`
	Votes     int    `json:"votes"`
}

type Client struct {
	Conn     *websocket.Conn
	Message  chan *Message
	ID       string `json:"id"`
	ChatID   int64  `json:"chat_id"`
	Username string `json:"username"`
}

type Chat struct {
	sync.Mutex
	Clients    map[string]*Client `json:"clients"`
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewChat() *Chat {
	return &Chat{
		Clients:    make(map[string]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

type Service interface {
	ReadMessage(ctx context.Context, chat *Chat, cl *Client)
	WriteMessage(ctx context.Context, cl *Client)
	IncVotes(ctx context.Context, id int64, chat *Chat)
	DecVotes(ctx context.Context, id int64, chat *Chat)
}

func (c *Chat) Run() {
	for {
		select {

		case client := <-c.Register:
			if _, ok := c.Clients[client.ID]; !ok {
				c.Clients[client.ID] = client
			}
		case client := <-c.Unregister:
			if _, ok := c.Clients[client.ID]; ok {
				if len(c.Clients) > 0 {
					c.Broadcast <- &Message{
						Content:  "User has left the chat",
						ChatID:   client.ChatID,
						Username: client.Username,
					}
				}
				delete(c.Clients, client.ID)
				close(client.Message)
			}
		case message := <-c.Broadcast:
			for _, cl := range c.Clients {
				cl.Message <- message
			}
		}
	}
}
