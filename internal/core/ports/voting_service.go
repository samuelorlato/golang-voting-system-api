package ports

import "github.com/gorilla/websocket"

type VotingService interface {
	HandleVote(*websocket.Conn)
	SetResultsConnections(*websocket.Conn)
}
