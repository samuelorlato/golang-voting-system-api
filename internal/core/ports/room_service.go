package ports

import (
	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/models"
)

type RoomService interface {
	CreateRoom(conn *websocket.Conn) string
	DeleteRoom(room *models.Room, index int)
	GetRoomIndex(conn *websocket.Conn, roomId string) *int
	GetRoom(conn *websocket.Conn, roomId string) *models.Room
	JoinAsVoter(conn *websocket.Conn, roomId string)
	JoinAsSpectator(conn *websocket.Conn, roomId string)
	Leave(conn *websocket.Conn, roomId string)
}
