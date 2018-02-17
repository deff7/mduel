package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"

	"sync"
	"time"
	//	"github.com/gorilla/websocket"
	//"net/url"
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

func readFromSocket(out chan<- []byte, quit <-chan struct{}, wg *sync.WaitGroup) {
	defer log.Println("Socket reader stopped")
	for {
		select {
		case <-quit:
			wg.Done()
			return
		default:
			out <- []byte("test")
			time.Sleep(time.Second)
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

	go readFromConsole(send, quit)
	wg.Add(1)
	go readFromSocket(receive, quit, &wg)

	for {
		select {
		case msg := <-send:
			log.Printf(" < %s", msg)
		case msg := <-receive:
			log.Printf(" > %s", msg)
		case <-interrupt:
			log.Println("Use Ctrl-D to stop")
		case <-quit:
			return
		}
	}
}
