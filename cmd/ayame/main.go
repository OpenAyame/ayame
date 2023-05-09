package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/OpenAyame/ayame"
	"golang.org/x/sync/errgroup"
)

const (
	ayameVersion = "2022.2.0"
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

	configFilePath := flag.String("c", "./config.ini", "ayame の設定ファイルへのパス(ini)")
	flag.Parse()

	config, err := ayame.NewConfig(*configFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// グローバルの logger に代入する
	if err := ayame.InitLogger(config); err != nil {
		log.Fatal(err)
	}

	server, err := ayame.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		addressAndPort := net.JoinHostPort(config.ListenPrometheusIPv4Address, strconv.Itoa(int(config.ListenPrometheusPortNumber)))
		return server.EchoPrometheus.Start(addressAndPort)
	})

	g.Go(func() error {
		return server.Start(ctx)
	})

	g.Go(func() error {
		return server.StartMatchServer(ctx)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}
