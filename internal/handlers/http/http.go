package httphandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
		hh.HandleVoteRequest(ctx.Writer, ctx.Request, roomId)
	})

	engine.GET("/spectate/*roomId", func(ctx *gin.Context) {
		roomId := ctx.Param("roomId")
		hh.HandleSpectateRequest(ctx.Writer, ctx.Request, roomId)
	})
}

func (hh *HTTPHandler) HandleVoteRequest(w http.ResponseWriter, r *http.Request, roomId string) {
	conn, err := hh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle
		return
	}

	if roomId == "/" {
		roomId = hh.roomService.CreateRoom(conn)
	} else {
		roomId = roomId[1:]
	}

	go hh.votingService.Vote(conn, roomId)
}

func (hh *HTTPHandler) HandleSpectateRequest(w http.ResponseWriter, r *http.Request, roomId string) {
	conn, err := hh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle
		return
	}

	if roomId == "/" {
		roomId = hh.roomService.CreateRoom(conn)
	} else {
		roomId = roomId[1:]
	}

	go hh.votingService.Spectate(conn, roomId)
}
