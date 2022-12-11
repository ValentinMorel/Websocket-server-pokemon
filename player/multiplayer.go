package player

import "time"

const (
	fps = 60
	// frameMs = 1000 / fps // Duration of a frame in milli seconds
	frameNs = int64(1e9) / fps // Duration of a frame in nano seconds
)

// A gas holds the set of players (particles)
type Multiplayer struct {
	// numPlayers   int
	Players map[*Player]bool
	Ids     map[string]*Player

	Remove chan *Player
}

func NewMultiplayer() *Multiplayer {
	self := &Multiplayer{}
	self.Players = map[*Player]bool{}
	self.Ids = map[string]*Player{}

	self.Remove = make(chan *Player)

	go self.Run()

	return self
}

func (self *Multiplayer) Run() {

	var start, timeElapsed, sleep int64

	for {
		start = time.Now().UnixNano()
		// We update the position of the players and we send it to them

		for p := range self.Players {
			p.update()
		}

		// Remove players
		removing := true
		for removing {
			select {
			case p := <-self.Remove:
				self.RemovePlayer(p)
			default:
				removing = false
			}
		}

		timeElapsed = time.Now().UnixNano() - start
		sleep = frameNs - timeElapsed

		// Necessary to run at a constant rate of fps
		time.Sleep(time.Duration(sleep) * time.Nanosecond)
		/*
		 time.Sleep takes a duration for an argument. time.Nanosecond
		   is of type Duration
		*/

	}
}

func (self *Multiplayer) Broadcast(event []byte) {
	for player, _ := range self.Players {
		if player.Ws != nil { // bots have nil
			player.send(event)
		}
	}
}

func (self *Multiplayer) AddPlayer(p *Player) {
	//var event []byte

	id := randString(8)

	self.Players[p] = true
	self.Ids[id] = p

	p.Id = id
	p.Others = self
	/*
		// We announce the new player
		event = CreatePlayerEvent(p)
		self.Broadcast(event)

		// To the new player we announce existing players
		for otherId, other := range self.Ids {
			if otherId != id {
				event = CreatePlayerEvent(other)
				p.send(event)
			}
		}
	*/
}

func (self *Multiplayer) RemovePlayer(p *Player) {
	delete(self.Players, p)
	delete(self.Ids, p.Id)
	// close(p.send)
	p.Others = nil
}

func (self *Multiplayer) NewId() string {
	var id string
	for {
		id = randString(8)
		if _, ok := self.Ids[id]; ok == false {
			// We only take unused ids
			return id
		}
	}
}
