package main

import (
	"github.com/gorilla/websocket"
	"log"
	"net/url"

	"bufio"

	"os"
	"os/signal"

	"sync"
)

func readFromTerminal(c chan<- []byte, quit <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("Terminal reader closed")
	scanner := bufio.NewScanner(os.Stdin)
	read := make(chan []byte)
	closeChan := make(chan struct{})
	go func() {
		defer close(closeChan)
		defer log.Println("Scanner thread closed")
		for {
			if scanner.Scan() {
				select {
				case <-quit:
					return
				case read <- scanner.Bytes():
				}
			} else {
				if err := scanner.Err(); err != nil {
					log.Println("Scanner:", err)
					return
				}
			}
		}
	}()

	for {
		select {
		case msg := <-read:
			c <- msg
		case <-closeChan:
			return
		}
	}
}

func main() {
	var wg sync.WaitGroup
	defer wg.Wait()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/socket"}
	log.Printf("Connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("Dial:", err)
	}
	defer c.Close()

	quit := make(chan struct{})
	send := make(chan []byte)

	go func() {
		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				log.Println("Read:", err)
				return
			}
			log.Printf(" > %s", msg)
		}
	}()

	wg.Add(1)
	go readFromTerminal(send, quit, &wg)

	for {
		select {
		case msg := <-send:
			err := c.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Println("Write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupted")
			err := c.WriteMessage(
				websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				log.Println("Write close:", err)
				return
			}
			close(quit)
			close(send)
			c.Close()
			return
		}
	}
}
