package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-voting-system-api/internal/core/services"
	httphandler "github.com/samuelorlato/golang-voting-system-api/internal/handlers/http"
)

func main() {
	engine := gin.Default()

	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	roomService := services.NewRoomService()
	votingService := services.NewVotingService(roomService)
	HTTPHandler := httphandler.NewHTTPHandler(&upgrader, roomService, votingService)
	HTTPHandler.SetRoutes(engine)

	engine.Run()
}
