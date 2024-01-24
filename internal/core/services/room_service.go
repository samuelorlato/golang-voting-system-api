package services

import (
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

func (r *roomService) CreateRoom() string {
	room := models.NewRoom()

	r.Lock()
	r.rooms = append(r.rooms, room)
	r.Unlock()

	return room.Id
}

func (r *roomService) DeleteRoom(conn *websocket.Conn, roomIndex *int) {
	room := r.GetRoom(roomIndex)
	if room == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
		return
	}

	r.Lock()
	r.rooms[*roomIndex] = r.rooms[len(r.rooms)-1]
	r.rooms = r.rooms[:len(r.rooms)-1]
	r.Unlock()
}

func (r *roomService) GetRoomIndex(roomId string) *int {
	for index, room := range r.rooms {
		if room.Id == roomId {
			return &index
		}
	}

	return nil
}

func (r *roomService) GetRoom(roomIndex *int) *models.Room {
	if roomIndex == nil {
		return nil
	}

	if *roomIndex > len(r.rooms) || *roomIndex < 0 {
		return nil
	}

	return r.rooms[*roomIndex]
}

func (r *roomService) JoinAsVoter(conn *websocket.Conn, roomIndex *int) {
	room := r.GetRoom(roomIndex)
	if room == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
		return
	}

	r.Lock()
	room.Voters[conn] = true
	r.Unlock()
}

func (r *roomService) JoinAsSpectator(conn *websocket.Conn, roomIndex *int) {
	room := r.GetRoom(roomIndex)
	if room == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
		return
	}

	r.Lock()
	room.Spectators[conn] = true
	r.Unlock()
}

func (r *roomService) Leave(conn *websocket.Conn, roomIndex *int) {
	room := r.GetRoom(roomIndex)
	if room == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
		return
	}

	r.Lock()
	delete(room.Spectators, conn)
	delete(room.Voters, conn)
	r.Unlock()

	conn.Close()

	if len(room.Voters) == 0 && len(room.Spectators) == 0 {
		r.DeleteRoom(conn, roomIndex)
	}
}