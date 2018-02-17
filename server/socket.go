package main

import (
	"github.com/deff7/mduel/server/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{}

func readFromSocket(c *websocket.Conn, out chan<- []byte) {
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("read: ", err)
			return
		}
		out <- msg
	}
}

func handleWebsocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade: ", err)
		return
	}
	defer c.Close()

	send := make(chan []byte)
	receive := make(chan []byte)

	state := game.Start(0, 1)

	go readFromSocket(c, receive)

	for {
		select {
		case msg := <-send:
			log.Println("s")
			err = c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("write: ", err)
				return
			}
		case msg := <-receive:
			log.Printf("rcv: %s", msg)
			go func() {
				state.HandleMessage(string(msg))
				send <- state.Encode()
			}()
		}
	}
}
