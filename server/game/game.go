package game

import (
	"encoding/json"
	"log"
)

type State struct {
	Players [2]*Player
}

func (s *State) Encode() []byte {
	res, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (s *State) HandleMessage(msg string) {
	p := s.Players[0]
	p.CheckSpelling(msg)
	if p.Cast() {
		s.Players[1].HP -= (p.powerLevel + 1) * 2
		p.powerLevel = 0
	}
}

func Start(id1, id2 int) *State {
	return &State{
		Players: [2]*Player{
			newPlayer(id1),
			newPlayer(id2),
		},
	}
}
