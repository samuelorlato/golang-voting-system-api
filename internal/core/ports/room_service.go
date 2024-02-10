package ports

import (
	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-voting-system-api/internal/core/models"
)

type RoomService interface {
	CreateRoom() string
	DeleteRoom(conn *websocket.Conn, roomIndex *int)
	GetRoomIndex(roomId string) *int
	GetRoom(roomIndex *int) *models.Room
	JoinAsVoter(conn *websocket.Conn, roomIndex *int)
	JoinAsSpectator(conn *websocket.Conn, roomIndex *int)
	Leave(conn *websocket.Conn, roomIndex *int)
}
