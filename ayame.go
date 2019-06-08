package main

import (
	"flag"
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"net/http"
)

var AyameVersion = "19.02.1"

type AyameOptions struct {
	LogDir         string
	LogName        string
	LogLevel       string
	Addr           string
	OverWsPingPong bool
}

var (
	// 起動時のオプション
	Options *AyameOptions
	logger  *logrus.Logger
)

// 初期化処理
func init() {
	logDir := flag.String("logDir", ".", "ayame log dir")
	logName := flag.String("logName", "ayame.log", "ayame log name")
	logLevel := flag.String("logLevel", "info", "ayame log name")
	overWsPingPong := flag.Bool("overWsPingPong", false, "with over-WS Ping-Pong")
	addr := flag.String("addr", "localhost:3000", " http service address")
	flag.Parse()
	Options = &AyameOptions{
		LogDir:         *logDir,
		LogName:        *logName,
		LogLevel:       *logLevel,
		Addr:           *addr,
		OverWsPingPong: *overWsPingPong,
	}
}

func main() {
	flag.Parse()
	args := flag.Args()
	logger = setupLogger()
	logger.Infof("WebRTC Signaling Server Ayame. version=%s", AyameVersion)
	logger.Infof("running on http://%s (Press Ctrl+C quit)", Options.Addr)
	// 引数の処理
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s", AyameVersion)
			return
		}
	}
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./sample/"+r.URL.Path[1:])
	})
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsHandler(hub, w, r)
	})
	logger.Fatal(http.ListenAndServe(Options.Addr, nil))
}
