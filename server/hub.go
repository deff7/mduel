package main

import (
	"errors"
	"log"
)

var maxRooms = 5

type Hub struct {
	rooms  []*Room
	lastID int

	register   chan *playerConnection
	unregister chan *playerConnection
}

func newHub() *Hub {
	return &Hub{
		register:   make(chan *playerConnection),
		unregister: make(chan *playerConnection),
	}
}

func (h *Hub) generatePlayerID() int {
	h.lastID += 1
	return h.lastID
}

func (h *Hub) findRoom() (*Room, error) {
	// find vacant room
	for _, r := range h.rooms {
		if r.vacant() {
			return r, nil
		}
	}
	// try to create new room
	if len(h.rooms) == maxRooms {
		return nil, errors.New("No empty rooms")
	}
	r := newRoom()
	h.rooms = append(h.rooms, r)
	return r, nil
}

func (h *Hub) registerConnection(c *playerConnection) {
	log.Println("registring connection...")
	room, err := h.findRoom()
	if err != nil {
		log.Println(err)
		return
	}

	c.id = h.generatePlayerID()
	room.add(c)
}

func (h *Hub) unregisterConnection(c *playerConnection) {
	r := c.room
	delete(r.playerConnections, c)
}

func (h *Hub) run() {
	for {
		select {
		case c := <-h.register:
			h.registerConnection(c)
		case c := <-h.unregister:
			h.unregisterConnection(c)
		}
	}
}
