package main

import (
	"github.com/deff7/mduel/server/game"
	"log"
	"time"
)

const (
	st_wait_for_players = iota
	st_game_started     = iota
)

type Room struct {
	state             int
	gameState         *game.State
	playerConnections map[*playerConnection]bool
	incomming         chan []byte
}

func (r *Room) vacant() bool {
	return len(r.playerConnections) < 2
}

func (r *Room) add(c *playerConnection) {
	r.playerConnections[c] = true
	c.room = r

	if !r.vacant() {
		r.state = st_game_started
		ids := []int{}
		for pconn := range r.playerConnections {
			ids = append(ids, pconn.id)
		}
		r.gameState = game.Start(ids[0], ids[1])
		log.Println("game started")
	}
}

func (r *Room) run() {
	log.Println("room running...")
	ticker := time.NewTicker(30 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if r.state == st_wait_for_players {
			} else if r.state == st_game_started {
				r.gameState.Update()
				for c := range r.playerConnections {
					go func() {
						c.send <- r.gameState.Encode()
					}()
				}
			}
		case msg := <-r.incomming:
			if r.state == st_game_started {
				go func() {
					log.Println("handle message from client")
					r.gameState.HandleMessage(string(msg))
				}()
			}
		}
	}
}

func newRoom() *Room {
	r := &Room{
		playerConnections: map[*playerConnection]bool{},
		incomming:         make(chan []byte),
		state:             st_wait_for_players,
	}
	go r.run()
	log.Println("new room created")
	return r
}
