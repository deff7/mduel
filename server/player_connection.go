package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

const (
	Connected     = iota
	WaitingRoom   = iota
	WaitingPlayer = iota
	StartingGame  = iota
	GameStarted   = iota
)

type playerConnection struct {
	id    int
	state int
	conn  *websocket.Conn
	room  *Room
	hub   *Hub
	send  chan []byte
}

var upgrader = websocket.Upgrader{}

func (p *playerConnection) read() {
	defer p.conn.Close()
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			p.hub.unregister <- p
			log.Println("read from socket:", err)
			return
		}
		p.room.incomming <- msg
	}
}

func (p *playerConnection) write() {
	defer p.conn.Close()
	for {
		select {
		case msg := <-p.send:
			err := p.conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				p.hub.unregister <- p
				log.Println("write to socket:", err)
				return
			}
		}
	}
}

func handleWebsocket(hub *Hub, w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}

	pconn := &playerConnection{
		conn:  c,
		send:  make(chan []byte),
		state: Connected,
		hub:   hub,
	}
	hub.register <- pconn

	go pconn.read()
	go pconn.write()
}
