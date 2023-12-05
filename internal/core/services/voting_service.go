package services

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/models"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
	"gopkg.in/validator.v2"
)

type votingService struct{}

func NewVotingService() ports.VotingService {
	return &votingService{}
}

var votes = make(map[int]int)
var clients = make(map[*websocket.Conn]bool)

var resultsConns []*websocket.Conn
var resultsConnMutex sync.Mutex

func (vs *votingService) HandleVote(conn *websocket.Conn) {
	clients[conn] = true

	go func() {
		for {
			var vote models.Vote
			err := conn.ReadJSON(&vote)
			if err != nil {
				delete(clients, conn)
				conn.Close()
				break
			}

			if len(resultsConns) > 0 {
				if err := validator.Validate(vote); err != nil {
					conn.WriteMessage(websocket.TextMessage, []byte("Invalid vote, probably invalid json format."))
				} else {
					votes[vote.Option]++
					broadcastVoteCountToResults()
				}
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Session is not open."))
			}
		}
	}()
}

func broadcastVoteCountToResults() {
	resultsConnMutex.Lock()
	defer resultsConnMutex.Unlock()

	for _, resultsConn := range resultsConns {
		if resultsConn == nil {
			return
		}

		for option, count := range votes {
			voteCount := models.VoteCount{Option: option, Count: count}
			err := resultsConn.WriteJSON(voteCount)
			if err != nil {
				delete(clients, resultsConn)
				resultsConn.Close()
			}
		}
	}
}

func (vs *votingService) SetResultsConnections(conn *websocket.Conn) {
	resultsConnMutex.Lock()
	defer resultsConnMutex.Unlock()
	resultsConns = append(resultsConns, conn)
}
