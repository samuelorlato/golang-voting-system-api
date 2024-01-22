package services

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/models"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
)

type roomService struct {
	rooms []*models.Room
	sync.RWMutex
}

func NewRoomService() ports.RoomService {
	return &roomService{
		rooms: []*models.Room{},
	}
}

func (r *roomService) CreateRoom(conn *websocket.Conn) string {
	room := models.NewRoom()

	r.Lock()
	r.rooms = append(r.rooms, room)
	r.Unlock()

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Room ID: %s", room.Id)))

	return room.Id
}

func (r *roomService) DeleteRoom(room *models.Room, index int) {
	r.Lock()
	r.rooms[index] = r.rooms[len(r.rooms)-1]
	r.rooms = r.rooms[:len(r.rooms)-1]
	r.Unlock()
}

func (r *roomService) GetRoomIndex(conn *websocket.Conn, roomId string) *int {
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

	return roomIndex
}

func (r *roomService) GetRoom(conn *websocket.Conn, roomId string) *models.Room {
	roomIndex := r.GetRoomIndex(conn, roomId)
	if roomIndex == nil {
		return nil
	}

	return r.rooms[*roomIndex]
}

func (r *roomService) JoinAsVoter(conn *websocket.Conn, roomId string) {
	roomIndex := r.GetRoomIndex(conn, roomId)
	if roomIndex == nil {
		return
	}

	r.Lock()
	r.rooms[*roomIndex].Voters[conn] = true
	r.Unlock()
}

func (r *roomService) JoinAsSpectator(conn *websocket.Conn, roomId string) {
	roomIndex := r.GetRoomIndex(conn, roomId)
	if roomIndex == nil {
		return
	}

	r.Lock()
	r.rooms[*roomIndex].Spectators[conn] = true
	r.Unlock()
}

func (r *roomService) Leave(conn *websocket.Conn, roomId string) {
	var foundIndex int
	var foundRoom *models.Room

	for index, room := range r.rooms {
		if room.Id == roomId {
			foundIndex = index
			foundRoom = room

			if room.Voters[conn] {
				r.Lock()
				delete(room.Voters, conn)
				r.Unlock()

				conn.Close()

				break
			} 
			
			if room.Spectators[conn] {
				r.Lock()
				delete(room.Spectators, conn)
				r.Unlock()

				conn.Close()

				break
			}
		}
	}

	if len(foundRoom.Voters) == 0 && len(foundRoom.Spectators) == 0 {
		r.DeleteRoom(foundRoom, foundIndex)
	}
}
