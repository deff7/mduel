package game

import (
	"encoding/json"
	"log"
)

var spellTypes = []string{"attack", "shield", "heal", "stun"}

type State struct {
	Players [2]*Player
}

type Message struct {
	PlayerID int
	Input    string
}

func (s *State) Encode() []byte {
	res, err := json.Marshal(s)
	if err != nil {
		log.Fatal(err)
	}
	return res
}

func (s *State) Update() {
	for _, p := range s.Players {
		enemy := p.enemy
		spell := p.Spell
		if spell.BoltSpeed > 0 {
			spell.Distance -= spell.BoltSpeed
			if spell.Distance <= 0 {
				enemy.Hurt((spell.powerLevel + 1) * 2)
				spell.Discharge()
			}
		}
	}
}

func (s *State) getPlayerByID(id int) *Player {
	for _, p := range s.Players {
		if p.ID == id {
			return p
		}
	}
	return nil
}

func (s *State) HandleMessage(msg []byte) {
	data := Message{}
	err := json.Unmarshal(msg, &data)
	if err != nil {
		log.Fatal(err)
	}

	p := s.getPlayerByID(data.PlayerID)
	spell := p.Spell
	spell.Check(data.Input)
	if spell.Cast() {
		log.Println("cast!")
		spell.BoltSpeed = 16 + spell.powerLevel*5
		spell.Distance = initialDistance
		p.updateSuggestions()
	}
}

func Start(id1, id2 int) *State {
	p1 := newPlayer(id1)
	p2 := newPlayer(id2)
	p1.enemy, p2.enemy = p2, p1
	return &State{Players: [2]*Player{p1, p2}}
}
