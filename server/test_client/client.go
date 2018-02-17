package main

import (
	"bufio"
	"log"
	"os"
	"os/signal"
)

func readFromConsole(out chan<- []byte, quit chan<- struct{}) {
	defer log.Println("Close scanner")
	defer close(quit)

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		out <- scanner.Bytes()
	}
	if err := scanner.Err(); err != nil {
		log.Println("Scanner:", err)
	}
}

func main() {
	defer log.Println("Bye bye")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	send := make(chan []byte)
	quit := make(chan struct{})

	go readFromConsole(send, quit)

	for {
		select {
		case msg := <-send:
			log.Printf(" < %s", msg)
		case <-interrupt:
			log.Println("Use Ctrl-D to stop")
		case <-quit:
			return
		}
	}
}
