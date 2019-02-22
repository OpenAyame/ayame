package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:3000", "http service address")

func main() {
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./sample/"+r.URL.Path[1:])
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(hub, w, r)
	})
	flag.Parse()
	log.SetFlags(0)
	log.Print(http.ListenAndServe(*addr, nil))
}
