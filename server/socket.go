package main

import (
	"github.com/deff7/mduel/server/game"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
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
	ticker := time.NewTicker(30 * time.Millisecond)
	for {
		select {
		//case msg := <-send:
		case <-ticker.C:
			state.Update()
			err = c.WriteMessage(websocket.TextMessage, state.Encode())
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
