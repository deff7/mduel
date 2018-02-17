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

func (s *State) Update() {
	enemy := s.Players[1]
	spell := s.Players[0].Spell
	if spell.BoltSpeed > 0 {
		spell.Distance -= spell.BoltSpeed
		if spell.Distance <= 0 {
			enemy.Hurt((spell.powerLevel + 1) * 2)
			spell.Discharge()
		}
	}
}

func (s *State) HandleMessage(msg string) {
	log.Println(msg)
	p := s.Players[0]

	spell := p.Spell
	spell.Check(msg)
	if spell.Cast() {
		log.Println("cast!")
		spell.BoltSpeed = 16 + spell.powerLevel*5
		spell.Distance = initialDistance
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
