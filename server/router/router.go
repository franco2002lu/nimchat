package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"server/internal/user"
	"server/internal/websocket"
	"time"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, wsHandler *websocket.Handler) {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/register", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.GET("/ws/joinChat/:chatId", wsHandler.JoinChat)
	r.GET("/ws/getClients/:chatId", wsHandler.GetClients)
	r.GET("/ws/:chatId/message/:messageId/upvote", wsHandler.IncVotes)
	r.GET("/ws/:chatId/message/:messageId/downvote", wsHandler.DecVotes)
}

func Start(addr string) error {
	return r.Run(addr)
}
