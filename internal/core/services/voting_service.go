package services

import (
	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/models"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
	"gopkg.in/validator.v2"
)

type votingService struct {
	roomService ports.RoomService
}

func NewVotingService(roomService ports.RoomService) ports.VotingService {
	return &votingService{
		roomService: roomService,
	}
}

func (vs *votingService) HandleVote(conn *websocket.Conn, roomId string) {
	vs.roomService.JoinAsVoter(conn, roomId)

	go func() {
		for {
			var vote models.Vote
			err := conn.ReadJSON(&vote)
			if err != nil {
				vs.roomService.Leave(conn, roomId)
			}

			if err := validator.Validate(vote); err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte("Invalid vote, probably invalid json format."))
			} else {
				vs.roomService.GetRoom(conn, roomId).Votes[vote.Option]++
				vs.broadcastVoteCountToResults(conn, roomId)
			}
		}
	}()
}

func (vs *votingService) broadcastVoteCountToResults(conn *websocket.Conn, roomId string) {
	// vs.m.Lock()
	// defer vs.m.Unlock()

	room := vs.roomService.GetRoom(conn, roomId)

	for conn := range room.Spectators {
		for option, count := range room.Votes {
			voteCount := models.VoteCount{Option: option, Count: count}
			err := conn.WriteJSON(voteCount)
			if err != nil {
				vs.roomService.Leave(conn, roomId)
			}
		}
	}
}

func (vs *votingService) SetResultsConnections(conn *websocket.Conn, roomId string) {
	// vs.m.Lock()
	// defer vs.m.Unlock()

	vs.roomService.JoinAsSpectator(conn, roomId)
}
