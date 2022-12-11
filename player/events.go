package player

import (
	"fmt"
)

// The events that we write to the clients
type wsEvent struct {
	//The json package only accesses the exported fields of struct types
	Event    string `json:"event"`
	Id       string `json:"id"`
	MapId    int    `json:"mapId"`
	X        int32  `json:"x"`
	Y        int32  `json:"y"`
	SpriteId int32  `json:"spriteId"`
}

func CreatePlayerEvent(p *Player) []byte {
	s := fmt.Sprintf(`{"event":"createPlayer", "id": "%v"}`,
		p.Id)
	return []byte(s)
}

func DestroyPlayerEvent(p *Player) []byte {
	s := fmt.Sprintf(`{"event":"destroyPlayer", "id":"%s"}`, p.Id)
	return []byte(s)
}

func UpdateEvent(p *Player) []byte {
	s := fmt.Sprintf(`{"event":"update", "id": "%v", "mapId": %v, "x": %v, "y": %v, "spriteId": %v}`, p.Id, p.MapId, p.X, p.Y, p.SpriteId)
	return []byte(s)
}
