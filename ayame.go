package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	logrus "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

const (
	// timeout は暫定的に 10 sec
	readHeaderTimeout = 10 * time.Second
)

var ayameVersion = "19.08.0"

type ayameOptions struct {
	LogDir         string `yaml:"log_dir"`
	LogName        string `yaml:"log_name"`
	LogLevel       string `yaml:"log_level"`
	Addr           string `yaml:"addr"`
	Port           int    `yaml:"port"`
	AuthWebhookURL string `yaml:"auth_webhook_url"`
	AllowOrigin    string `yaml:"allow_origin"`
}

var (
	options *ayameOptions
	logger  *logrus.Logger = newLogger()
)

// 初期化処理
func init() {
	configFilePath := flag.String("c", "./config.yaml", "ayame の設定ファイルへのパス(yaml)")
	flag.Parse()
	// yaml ファイルを読み込み
	buf, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		// 読み込めない場合 Fatal で終了
		log.Fatal("cannot open config file, err=", err)
	}
	// yaml をパース

	if err := yaml.Unmarshal(buf, &options); err != nil {
		// パースに失敗した場合 Fatal で終了
		log.Fatal("cannot parse config file, err=", err)
	}

	if err := setupLogger(logger); err != nil {
		log.Fatal(err)
	}
}

func main() {
	args := flag.Args()
	// 引数の処理
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s", ayameVersion)
			return
		}
	}

	url := fmt.Sprintf("%s:%d", options.Addr, options.Port)
	logger.Infof("WebRTC Signaling Server Ayame. version=%s", ayameVersion)
	logger.Infof("running on http://%s (Press Ctrl+C quit)", url)
	hub := newHub()
	go hub.run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./sample/"+r.URL.Path[1:])
	})
	// /ws エンドポイントは将来的に /signaling に統一するが、互換性のために残しておく
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		signalingHandler(hub, w, r)
	})
	http.HandleFunc("/signaling", func(w http.ResponseWriter, r *http.Request) {
		signalingHandler(hub, w, r)
	})
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: readHeaderTimeout}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal(err)
	}
}
