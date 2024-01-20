package services

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/models"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
)

type roomService struct {
	rooms []*models.Room
}

func NewRoomService() ports.RoomService {
	return &roomService{
		rooms: []*models.Room{},
	}
}

func (r *roomService) CreateRoom(conn *websocket.Conn) string {
	room := models.NewRoom()
	r.rooms = append(r.rooms, room)

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Room ID: %s", room.Id)))

	return room.Id
}

func (r *roomService) DeleteRoom(room *models.Room, index int) {
	r.rooms[index] = r.rooms[len(r.rooms)-1]
	r.rooms = r.rooms[:len(r.rooms)-1]
}

func (r *roomService) GetRoomIndex(conn *websocket.Conn, roomId string) int {
	var roomIndex *int

	for index, room := range r.rooms {
		if room.Id == roomId {
			roomIndex = &index
			break
		}
	}

	if roomIndex == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
	}

	return *roomIndex
}

func (r *roomService) GetRoom(conn *websocket.Conn, roomId string) *models.Room {
	var roomIndex *int

	for index, room := range r.rooms {
		if room.Id == roomId {
			roomIndex = &index
			break
		}
	}

	if roomIndex == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
	}

	return r.rooms[*roomIndex]
}

func (r *roomService) JoinAsVoter(conn *websocket.Conn, roomId string) {
	roomIndex := r.GetRoomIndex(conn, roomId)
	r.rooms[roomIndex].Voters[conn] = true
}

func (r *roomService) JoinAsSpectator(conn *websocket.Conn, roomId string) {
	roomIndex := r.GetRoomIndex(conn, roomId)
	r.rooms[roomIndex].Spectators[conn] = true
}

func (r *roomService) Leave(conn *websocket.Conn, roomId string) {
	for index, room := range r.rooms {
		if room.Id == roomId {
			for k := range room.Voters {
				if k == conn {
					delete(room.Voters, conn)
					conn.Close()
					break
				}
			}

			for k := range room.Spectators {
				if k == conn {
					delete(room.Spectators, conn)
					conn.Close()
					break
				}
			}

			if len(room.Voters) == 0 && len(room.Spectators) == 0 {
				r.DeleteRoom(room, index)
			}
		}
	}
}
