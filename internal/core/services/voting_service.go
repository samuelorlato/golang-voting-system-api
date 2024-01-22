package services

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/models"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
	"gopkg.in/validator.v2"
)

type votingService struct {
	roomService ports.RoomService
	sync.RWMutex
}

func NewVotingService(roomService ports.RoomService) ports.VotingService {
	return &votingService{
		roomService: roomService,
	}
}

func (vs *votingService) Vote(conn *websocket.Conn, roomId string) {
	if vs.roomService.GetRoom(conn, roomId) == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
		return
	}

	vs.roomService.JoinAsVoter(conn, roomId)

	for {
		var vote models.Vote
		err := conn.ReadJSON(&vote)
		if err != nil {
			vs.roomService.Leave(conn, roomId)
			break
		}

		if err := validator.Validate(vote); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid vote, probably invalid json format."))
		} else {
			room := vs.roomService.GetRoom(conn, roomId)

			vs.addVote(room, vote)
			vs.broadcastVoteCountToSpectators(conn, room)
		}
	}
}

func (vs *votingService) addVote(room *models.Room, vote models.Vote) {
	vs.Lock()
	room.Votes[vote.Option]++
	vs.Unlock()
}

func (vs *votingService) broadcastVoteCountToSpectators(conn *websocket.Conn, room *models.Room) {
	vs.Lock()
	for spectatorConn := range room.Spectators {
		for option, count := range room.Votes {
			voteCount := models.VoteCount{Option: option, Count: count}

			err := spectatorConn.WriteJSON(voteCount)
			if err != nil {
				vs.roomService.Leave(spectatorConn, room.Id)
				break
			}
		}
	}
	vs.Unlock()
}

func (vs *votingService) Spectate(conn *websocket.Conn, roomId string) {
	if vs.roomService.GetRoom(conn, roomId) == nil {
		conn.WriteMessage(websocket.TextMessage, []byte("Room not found."))
		conn.Close()
		return
	}

	vs.roomService.JoinAsSpectator(conn, roomId)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			vs.roomService.Leave(conn, roomId)
			break
		}
	}
}
