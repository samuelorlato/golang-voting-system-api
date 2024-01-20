package ports

import (
	"github.com/gorilla/websocket"
)

type VotingService interface {
	HandleVote(conn *websocket.Conn, roomId string)
	SetResultsConnections(conn *websocket.Conn, roomId string)
}
