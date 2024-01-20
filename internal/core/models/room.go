package models

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Room struct {
	Id         string
	Votes      map[int]int
	Voters     map[*websocket.Conn]bool
	Spectators map[*websocket.Conn]bool
}

func NewRoom() *Room {
	room := Room{
		Id:         uuid.NewString(),
		Votes:      map[int]int{},
		Voters:     map[*websocket.Conn]bool{},
		Spectators: map[*websocket.Conn]bool{},
	}

	return &room
}
