package voting

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Vote struct {
	Option int `json:"option"`
}

type VoteCount struct {
	Option int `json:"option"`
	Count  int `json:"count"`
}

var votes = make(map[int]int)
var clients = make(map[*websocket.Conn]bool)

var resultsConns []*websocket.Conn
var resultsConnMutex sync.Mutex

func HandleVote(conn *websocket.Conn) {
	clients[conn] = true

	go func() {
		for {
			var vote Vote
			err := conn.ReadJSON(&vote)
			if err != nil {
				delete(clients, conn)
				conn.Close()
				break
			}

			if len(resultsConns) > 0 {
				votes[vote.Option]++
				broadcastVoteCountToResults()
			} else {
				conn.WriteMessage(websocket.TextMessage, []byte("Session is not open."))
			}
		}
	}()
}

func SetResultsConnections(conn *websocket.Conn) {
	resultsConnMutex.Lock()
	defer resultsConnMutex.Unlock()
	resultsConns = append(resultsConns, conn)
}

func broadcastVoteCountToResults() {
	resultsConnMutex.Lock()
	defer resultsConnMutex.Unlock()

	for _, resultsConn := range resultsConns {
		if resultsConn == nil {
			return
		}

		for option, count := range votes {
			voteCount := VoteCount{Option: option, Count: count}
			err := resultsConn.WriteJSON(voteCount)
			if err != nil {
				delete(clients, resultsConn)
				resultsConn.Close()
			}
		}
	}
}
