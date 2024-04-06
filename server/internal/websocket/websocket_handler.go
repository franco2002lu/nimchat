package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"server/db"
	"strconv"
)

type Handler struct {
	chat *Chat
}

func NewHandler(c *Chat) *Handler {
	return &Handler{
		chat: c,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// orgin := r. Header. Get ("Origin")
		// return origin == "http://localhost:3000"
		return true //todo: replace with localhost:3000 or something later,check second video
	},
}

func (h *Handler) JoinChat(c *gin.Context) {
	connection, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	chatID := c.Param("chatID")
	clientID := uuid.New().String()
	username := c.Query("username")

	NewChatID, err := strconv.ParseInt(chatID, 10, 64)
	if err != nil {
		NewChatID = 1
	}

	cl := &Client{
		Conn:     connection,
		Message:  make(chan *Message, 10),
		ID:       clientID,
		ChatID:   NewChatID,
		Username: username,
	}

	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatalf("Database connection failed: %s", err)
	}

	wsRepo := NewRepo(dbConnection.GetDB())
	wsService := NewService(wsRepo)

	h.chat.Register <- cl

	go wsService.WriteMessage(c, cl)
	wsService.ReadMessage(c, h.chat, cl)
}

type ClientResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	clients := make([]ClientResponse, 0)

	for _, cl := range h.chat.Clients {
		clients = append(clients, ClientResponse{
			ID:       cl.ID,
			Username: cl.Username,
		})
	}

	c.JSON(http.StatusOK, clients)
}

func (h *Handler) IncVotes(c *gin.Context) {
	// mutex to handle race conditions
	h.chat.Lock()
	defer h.chat.Unlock()

	messageId := c.Param("messageId")

	_, err := strconv.Atoi(messageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: Invalid message id": messageId})
		return
	}
	intMsgId, err := strconv.ParseInt(messageId, 10, 64)
	if err != nil {
		log.Println("Conversion error:", err)
	}

	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatalf("Database connection failed: %s", err)
	}

	wsRepo := NewRepo(dbConnection.GetDB())
	wsService := NewService(wsRepo)

	wsService.IncVotes(c, intMsgId, h.chat)
}

func (h *Handler) DecVotes(c *gin.Context) {
	// mutex to handle race conditions
	h.chat.Lock()
	defer h.chat.Unlock()

	messageId := c.Param("messageId")

	_, err := strconv.Atoi(messageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error: Invalid message id": messageId})
		return
	}
	intMsgId, err := strconv.ParseInt(messageId, 10, 64)
	if err != nil {
		log.Println("Conversion error:", err)
	}

	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatalf("Database connection failed: %s", err)
	}

	wsRepo := NewRepo(dbConnection.GetDB())
	wsService := NewService(wsRepo)

	wsService.DecVotes(c, intMsgId, h.chat)
}
