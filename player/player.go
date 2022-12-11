package player

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/websocket"
)

const letters = "abcdefghijklmnopqrstuvwxyz"

type Player struct {
	Ws       *websocket.Conn
	Output   chan []byte
	Id       string `json:"id"`
	MapId    int    `json:"mapId"`
	X        int32  `json:"x"`
	Y        int32  `json:"y"`
	SpriteId int32  `json:"spriteId"`
	*wsEvent
	Others *Multiplayer
}

func NewPlayer(ws *websocket.Conn) *Player {
	self := &Player{}
	self.Ws = ws

	if ws != nil {
		self.Output = make(chan []byte, 64) // buffered so it doesn't block
	} else {
		self.Output = nil
	}
	return self
}

func (self *Player) send(event []byte) {
	if self.Output == nil {
		return
	}

	self.Output <- event
}

func (self *Player) Reader() {
	for {
		_, event, err := self.Ws.ReadMessage()
		if err != nil {
			break
		}
		json.Unmarshal(event, self)
		//log.Printf("%s -> %s\n", self.Ws.RemoteAddr(), event)
		/* go handleWsEvent(c, j)
		No longer necessary, tick uses directly the new sate */
	}
	self.Ws.Close()
	// We need to remove a player more carefully!
}
func (self *Player) Writer() {
	for event := range self.Output {
		w, _ := self.Ws.NextWriter(websocket.TextMessage)

		w.Write(event)

		//log.Printf("%s <- %s\n", self.Ws.RemoteAddr(), event)
	}
	self.Ws.Close()
}

// Used to generate random id
func randString(n int) string {
	b := make([]byte, n)
	// The result of letters[i] isof type uint8 (byte)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
func (self *Player) update() {
	event := UpdateEvent(self)

	for p := range self.Others.Players {
		if p.Ws != nil && p.Id != self.Id { // Not sent to bots
			if self.isNear(p) {
				p.send(event)
			}
		}
	}
}

func (self *Player) isNear(other *Player) bool {

	if self.MapId == other.MapId {
		return true
	} else {
		return false
	}

	return true
}
