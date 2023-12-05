package httphandler

import (
	"net/http"

	"github.com/samuelorlato/golang-electoral-system-api/internal/core/ports"
)

type HTTPHandler struct {
	upgrader      ports.HTTPToSocketConnectionUpgrader
	votingService ports.VotingService
}

func NewHTTPHandler(upgrader ports.HTTPToSocketConnectionUpgrader, votingService ports.VotingService) *HTTPHandler {
	return &HTTPHandler{
		upgrader:      upgrader,
		votingService: votingService,
	}
}

func (hh *HTTPHandler) HandleVoteRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := hh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle
		return
	}

	hh.votingService.HandleVote(conn)
}

func (hh *HTTPHandler) HandleResultsRequest(w http.ResponseWriter, r *http.Request) {
	conn, err := hh.upgrader.Upgrade(w, r, nil)
	if err != nil {
		// handle
		return
	}

	hh.votingService.SetResultsConnections(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			// handle
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			// handle
			return
		}
	}
}
