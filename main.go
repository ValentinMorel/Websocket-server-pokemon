package main

import (
	"gameboy-server/player"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = &websocket.Upgrader{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		EnableCompression: true,
	}
)

var (
	// Declared in top level so that it is visible for the whole package
	mainMulti *player.Multiplayer
)

func init() {
	// Declared in top level so that it is visible for the whole package
	mainMulti = player.NewMultiplayer()
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	// Handler to /ws
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	addr := ws.RemoteAddr()

	log.Printf("Websocket accepted: %s\n", addr)
	newPlayer := player.NewPlayer(ws)
	mainMulti.AddPlayer(newPlayer)
	ws.WriteMessage(websocket.TextMessage, player.CreatePlayerEvent(newPlayer))
	go newPlayer.Writer()
	go newPlayer.Reader()

	// Once reader returns the connection is finalized
	log.Printf("Websocket finalized: %s\n", addr)

}

func RunServer() {

	http.HandleFunc("/socket", wsHandler)
	log.Fatal(http.ListenAndServe("localhost:8081", nil))
}

func main() {
	RunServer()
}
