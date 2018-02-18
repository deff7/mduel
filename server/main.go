package main

import (
	"html/template"
	"log"
	"net/http"
)

func home(rw http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("data/templates/main.html")
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	t.Execute(rw, nil)
}

func main() {
	http.HandleFunc("/", home)

	fileServer := http.StripPrefix("/static/", http.FileServer(http.Dir("data/static")))
	http.Handle("/static/", fileServer)

	hub := newHub()
	http.HandleFunc("/socket", func(w http.ResponseWriter, r *http.Request) {
		handleWebsocket(hub, w, r)
	})

	go hub.run()
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
