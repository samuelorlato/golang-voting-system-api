package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id         string
	TotalVotes int
	Votes      map[string]int
	Voters     map[*websocket.Conn]bool
	Spectators map[*websocket.Conn]bool
}

func NewRoom() *Room {
	room := Room{
		Id:         uuid.NewString(),
		TotalVotes: 0,
		Votes:      map[string]int{},
		Voters:     map[*websocket.Conn]bool{},
		Spectators: map[*websocket.Conn]bool{},
	}

	return &room
}
