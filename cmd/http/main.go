package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/internal/core/services"
	httphandler "github.com/samuelorlato/golang-electoral-system-api/internal/handlers/http"
)

func main() {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	votingService := services.NewVotingService()
	HTTPHandler := httphandler.NewHTTPHandler(&upgrader, votingService)

	http.HandleFunc("/vote", HTTPHandler.HandleVoteRequest)
	http.HandleFunc("/results", HTTPHandler.HandleResultsRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
