package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"

	"github.com/gorilla/websocket"
	"net/url"
	"sync"

	"encoding/json"
	tm "github.com/buger/goterm"
	//	"github.com/davecgh/go-spew/spew"
	"github.com/deff7/mduel/server/game"
	//	"time"
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

	gameState := game.State{}
	lastWord := ""

	for {
		select {
		case msg := <-send:
			lastWord = string(msg)
			c.WriteMessage(websocket.TextMessage, msg)
		case msg := <-receive:
			err := json.Unmarshal(msg, &gameState)
			if err != nil {
				log.Fatal(err)
			}
			tm.Clear()
			tm.MoveCursor(1, 1)
			tm.Printf("Enemy HP: %d", gameState.Players[1].HP)
			spell := gameState.Players[0].Spell
			if word := spell.NextWord; word != "" {
				tm.Printf("\nNext word of power: %s", word)
			}
			if spell.BoltSpeed > 0 {
				tm.Printf("\nBolt distance: %d", spell.Distance)
			}
			if lastWord != "" {
				tm.Printf("\nLast word: %s", lastWord)
			}
			tm.Flush()
			//time.Sleep(100 * time.Millisecond)
		case <-interrupt:
			log.Println("Use Ctrl-D to stop")
		case <-quit:
			return
		}
	}
}
