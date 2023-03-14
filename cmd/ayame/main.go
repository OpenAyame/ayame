package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/OpenAyame/ayame"
	"golang.org/x/sync/errgroup"
	"gopkg.in/ini.v1"
)

const (
	ayameVersion = "2022.2.0"
	// timeout は暫定的に 10 sec
	readHeaderTimeout = 10 * time.Second
)

func main() {
	args := flag.Args()
	// 引数の処理
	if len(args) > 0 {
		if args[0] == "version" {
			fmt.Printf("WebRTC Signaling Server Ayame version %s", ayameVersion)
			return
		}
	}

	configFilePath := flag.String("c", "./ayame.ini", "ayame の設定ファイルへのパス(ini)")
	flag.Parse()

	iniConfig, err := ini.InsensitiveLoad(*configFilePath)
	if err != nil {
		log.Fatal("Cannot parse config file, err=", err)
	}

	config := new(ayame.Config)
	if err := ayame.InitConfig(*iniConfig, config); err != nil {
		log.Fatal(err)
	}

	// グローバルの logger に代入する
	if err := ayame.InitLogger(*config); err != nil {
		log.Fatal(err)
	}

	ayame.SetDefaultsConfig(*config)

	server, err := ayame.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return server.Start(ctx)
	})

	g.Go(func() error {
		return server.StartMatchServer()
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
