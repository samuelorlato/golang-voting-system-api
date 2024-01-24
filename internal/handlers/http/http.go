package httphandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
)

type HTTPHandler struct {
	upgrader      ports.HTTPToSocketConnectionUpgrader
	roomService   ports.RoomService
	votingService ports.VotingService
}

func NewHTTPHandler(upgrader ports.HTTPToSocketConnectionUpgrader, roomService ports.RoomService, votingService ports.VotingService) *HTTPHandler {
	return &HTTPHandler{
		upgrader:      upgrader,
		roomService:   roomService,
		votingService: votingService,
	}
}

func (hh *HTTPHandler) SetRoutes(engine *gin.Engine) {
	engine.GET("/vote/*roomId", func(ctx *gin.Context) {
		roomId := ctx.Param("roomId")
		hh.HandleVoteRequest(ctx, roomId)
	})

	engine.GET("/spectate/*roomId", func(ctx *gin.Context) {
		roomId := ctx.Param("roomId")
		hh.HandleSpectateRequest(ctx, roomId)
	})
}

func (hh *HTTPHandler) HandleVoteRequest(ctx *gin.Context, roomId string) {
	conn, err := hh.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// handle
		return
	}

	if roomId == "/" {
		roomId = hh.roomService.CreateRoom()
		conn.WriteMessage(websocket.TextMessage, []byte("Room ID: "+roomId))
	} else {
		roomId = roomId[1:]
	}

	go hh.votingService.Vote(conn, roomId)
}

func (hh *HTTPHandler) HandleSpectateRequest(ctx *gin.Context, roomId string) {
	conn, err := hh.upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		// handle
		return
	}

	if roomId == "/" {
		roomId = hh.roomService.CreateRoom()
		conn.WriteMessage(websocket.TextMessage, []byte("Room ID: "+roomId))
	} else {
		roomId = roomId[1:]
	}

	go hh.votingService.Spectate(conn, roomId)
}
