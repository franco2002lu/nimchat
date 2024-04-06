package main

import (
	"log"
	"server/db"
	"server/internal/user"
	"server/internal/websocket"
	"server/router"
)

func main() {
	dbConnection, err := db.NewDB()
	if err != nil {
		log.Fatalf("Database connection failed: %s", err)
	}

	usrRepo := user.NewRepo(dbConnection.GetDB())
	usrService := user.NewService(usrRepo)
	usrHandler := user.NewHandler(usrService)

	chat := websocket.NewChat()
	websocketHandler := websocket.NewHandler(chat)
	go chat.Run()

	router.InitRouter(usrHandler, websocketHandler)
	router.Start("0.0.0.0:8080")
}
