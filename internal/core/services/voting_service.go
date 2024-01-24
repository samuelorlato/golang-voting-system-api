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
	roomIndex := vs.roomService.GetRoomIndex(roomId)

	vs.roomService.JoinAsVoter(conn, roomIndex)

	for {
		var vote models.Vote
		err := conn.ReadJSON(&vote)
		if err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid vote: "+err.Error()))
			vs.roomService.Leave(conn, roomIndex)
			break
		}

		if err := validator.Validate(vote); err != nil {
			conn.WriteMessage(websocket.TextMessage, []byte("Invalid vote: "+err.Error()))
			vs.roomService.Leave(conn, roomIndex)
			break
		}

		room := vs.roomService.GetRoom(roomIndex)

		vs.addVote(room, vote)
		vs.broadcastVoteCountToSpectators(conn, room, roomIndex)
	}
}

func (vs *votingService) addVote(room *models.Room, vote models.Vote) {
	vs.Lock()
	room.Votes[vote.Option]++
	room.TotalVotes++
	vs.Unlock()
}

func (vs *votingService) broadcastVoteCountToSpectators(conn *websocket.Conn, room *models.Room, roomIndex *int) {
	vs.Lock()
	for spectatorConn := range room.Spectators {
		votesAndTotalVotesMap := map[string]interface{}{}
		votesAndTotalVotesMap["totalVotes"] = room.TotalVotes

		votesMap := map[string]models.VoteCount{}
		for option, count := range room.Votes {
			voteCount := models.VoteCount{Count: count, PercentageFromTotal: float64(count) / float64(room.TotalVotes)}
			votesMap[option] = voteCount
		}

		votesAndTotalVotesMap["separatedVotes"] = votesMap

		err := spectatorConn.WriteJSON(votesAndTotalVotesMap)
		if err != nil {
			vs.roomService.Leave(spectatorConn, roomIndex)
			break
		}
	}
	vs.Unlock()
}

func (vs *votingService) Spectate(conn *websocket.Conn, roomId string) {
	roomIndex := vs.roomService.GetRoomIndex(roomId)

	vs.roomService.JoinAsSpectator(conn, roomIndex)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			vs.roomService.Leave(conn, roomIndex)
			break
		}
	}
}
