package main

import (
	"encoding/json"
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
	msg, err := json.Marshal(struct {
		PlayerID    int
		PlayerIndex int
	}{
		PlayerID:    c.id,
		PlayerIndex: len(r.playerConnections) - 1,
	})
	if err != nil {
		log.Fatal(err)
	}
	c.send <- msg

	if !r.vacant() {
		r.state = st_game_started
		ids := []int{}
		for pconn := range r.playerConnections {
			msg, err := json.Marshal(struct {
				Ready bool
			}{
				Ready: true,
			})
			if err != nil {
				log.Fatal(err)
			}
			pconn.send <- msg

			ids = append(ids, pconn.id)
		}
		r.gameState = game.Start(ids[0], ids[1])
		log.Println("game started")
	}
}

func (r *Room) run() {
	log.Println("room running...")
	ticker := time.NewTicker(100 * time.Millisecond)
	for {
		select {
		case <-ticker.C:
			if r.state == st_wait_for_players {
			} else if r.state == st_game_started {
				r.gameState.Update()
				msg := r.gameState.Encode()
				for c := range r.playerConnections {
					c.send <- msg
				}
			}
		case msg := <-r.incomming:
			if r.state == st_game_started {
				go func() {
					log.Println("handle message from client")
					r.gameState.HandleMessage(msg)
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
