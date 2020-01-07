package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/yaml.v2"
)

const (
	ayameVersion = "2020.1"
	// timeout は暫定的に 10 sec
	readHeaderTimeout = 10 * time.Second
)

var (
	config          *ayameConfig
	logger          *zerolog.Logger
	signalingLogger *zerolog.Logger
	webhookLogger   *zerolog.Logger
)

// 初期化処理
func init() {
	testing.Init()
	configFilePath := flag.String("c", "./ayame.yaml", "ayame の設定ファイルへのパス(yaml)")
	flag.Parse()
	// yaml ファイルを読み込み
	buf, err := ioutil.ReadFile(*configFilePath)
	if err != nil {
		// 読み込めない場合 Fatal で終了
		log.Fatal("Cannot open config file, err=", err)
	}
	if err := yaml.Unmarshal(buf, &config); err != nil {
		log.Fatal("Cannot parse config file, err=", err)
	}

	setDefaultsConfig()

	// グローバルの logger に代入する
	logger, err = initLogger()
	if err != nil {
		log.Fatal(err)
	}

	// グローバルの signalingLogger に代入する
	signalingLogger, err = initSignalingLogger()
	if err != nil {
		log.Fatal(err)
	}

	if config.AuthnWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.AuthnWebhookURL); err != nil {
			log.Fatal(err)
		}
	}

	if config.DisconnectWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.DisconnectWebhookURL); err != nil {
			log.Fatal(err)
		}
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

	// コンフィグのデフォルト値を追加する
	logger.Info().Str("log_dir", config.LogDir).Msg("AyameConf")
	logger.Info().Str("log_name", config.LogName).Msg("AyameConf")
	logger.Info().Str("log_level", config.LogLevel).Msg("AyameConf")
	logger.Info().Str("signaling_log_name", config.SignalingLogName).Msg("AyameConf")
	logger.Info().Str("listen_ipv4_address", config.ListenIPv4Address).Msg("AyameConf")
	logger.Info().Int("listen_port_number", config.ListenPortNumber).Msg("AyameConf")
	logger.Info().Str("authn_webhook_url", config.AuthnWebhookURL).Msg("AyameConf")
	logger.Info().Str("disconnect_webhook_url", config.DisconnectWebhookURL).Msg("AyameConf")
	logger.Info().Str("webhook_log_name", config.WebhookLogName).Msg("AyameConf")
	logger.Info().Int("webhook_request_timeout", config.WebhookRequestTimeout).Msg("AyameConf")
	logger.Info().Str("allow_origin", config.AllowOrigin).Msg("AyameConf")

	// URL の生成
	url := fmt.Sprintf("%s:%d", config.ListenIPv4Address, config.ListenPortNumber)

	go server()

	http.HandleFunc("/signaling", func(w http.ResponseWriter, r *http.Request) {
		signalingHandler(w, r)
	})
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: readHeaderTimeout}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Err(err)
	}
}
