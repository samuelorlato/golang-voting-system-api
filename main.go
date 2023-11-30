package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/samuelorlato/golang-electoral-system-api/core/voting"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc("/vote", vote)
	http.HandleFunc("/results", results)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func vote(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	voting.HandleVote(conn)
}

func results(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	voting.SetResultsConnections(conn)

	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
