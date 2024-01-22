package ports

import (
	"github.com/gorilla/websocket"
)

type VotingService interface {
	Vote(conn *websocket.Conn, roomId string)
	Spectate(conn *websocket.Conn, roomId string)
}
