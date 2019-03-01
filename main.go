package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = flag.String("addr", "localhost:3000", " http service address")
var AyameVersion = "19.02.1"

func main() {
	flag.Parse()
	args := flag.Args()
	// 引数の処理
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s", AyameVersion)
			return
		}
	}
	log.SetFlags(0)
	log.Printf("WebRTC Signaling Server Ayame\n version %s\n running on http://%s (Press Ctrl+C quit)\n", AyameVersion, *addr)
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./sample/"+r.URL.Path[1:])
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(*addr, nil))
}
