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

func main() {
	// bin/ayame -V
	showVersion := flag.Bool("V", false, "バージョン")

	// bin/ayame -C config.ini
	configFilePath := flag.String("C", "./config.ini", "設定ファイルへのパス")
	flag.Parse()

	if *showVersion {
		fmt.Printf("WebRTC Signaling Server Ayame version %s\n", ayame.Version)
		return
	}

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
