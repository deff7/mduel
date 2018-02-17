package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"net/url"
	"sync"
	//"time"
)

func readFromConsole(out chan<- []byte, quit chan<- struct{}) {
	defer log.Println("Console reading stopped")
	defer close(quit)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out <- scanner.Bytes()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Scanner:", err)
	}
}

func readFromSocket(conn *websocket.Conn, out chan<- []byte, quit <-chan struct{}, wg *sync.WaitGroup) {
	defer log.Println("Socket reader stopped")
	defer wg.Done()

	for {
		select {
		case <-quit:
			return
		default:
			_, msg, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			out <- msg
		}
	}
}

func main() {
	var wg sync.WaitGroup
	defer log.Println("Bye bye")
	defer wg.Wait()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	send := make(chan []byte)
	receive := make(chan []byte)
	quit := make(chan struct{})

	u := url.URL{Scheme: "ws", Host: "localhost:3000", Path: "/socket"}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		c.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		c.Close()
	}()

	go readFromConsole(send, quit)
	wg.Add(1)
	go readFromSocket(c, receive, quit, &wg)

	for {
		select {
		case msg := <-send:
			log.Printf(" < %s", msg)
			c.WriteMessage(websocket.TextMessage, msg)
		case msg := <-receive:
			log.Printf(" > %s", msg)
		case <-interrupt:
			log.Println("Use Ctrl-D to stop")
		case <-quit:
			return
		}
	}
}
