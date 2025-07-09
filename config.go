package ayame

import (
	_ "embed"
	"net/url"

	"log/slog"
	"gopkg.in/ini.v1"
)

//go:embed VERSION
var Version string

const (
	defaultLogDir  = "."
	defaultLogName = "ayame.jsonl"

	// 200 MB
	defaultLogRotateMaxSize = 200
	// 7 ファイル
	defaultLogRotateMaxBackups = 7
	// 30 日
	defaultLogRotateMaxAge = 30

	defaultSignalingLogName = "signaling.jsonl"

	defaultListenIPv4Address = "0.0.0.0"
	defaultListenPortNumber  = 3000

	defaultWebSocketReadTimeoutSec  = 90
	defaultWebSocketPongTimeoutSec  = 60
	defaultWebSocketPingIntervalSec = 5

	defaultWebhookRequestTimeout = 5
	defaultWebhookLogName        = "webhook.jsonl"

	defaultListenPrometheusIPv4Address = "0.0.0.0"
	defaultListenPrometheusPortNumber  = 4000
)

var defaultSignalingLogFilters = []string{
	"register",
	"offer",
	"answer",
	"candidate",
	"connected",
	"message",
}

type Config struct {
	Debug bool `ini:"debug"`

	LogDir              string `ini:"log_dir"`
	LogName             string `ini:"log_name"`
	LogStdout           bool   `ini:"log_stdout"`
	LogRotateMaxSize    int    `ini:"log_rotate_max_size"`
	LogRotateMaxBackups int    `ini:"log_rotate_max_backups"`
	LogRotateMaxAge     int    `ini:"log_rotate_max_age"`
	LogRotateCompress   bool   `ini:"log_rotate_compress"`

	LogMessageKeyName   string `ini:"log_message_key_name"`
	LogTimestampKeyName string `ini:"log_timestamp_key_name"`

	SignalingLogName    string   `ini:"signaling_log_name"`
	SignalingLogFilters []string `ini:"signaling_log_filters"`

	DebugConsoleLog     bool `ini:"debug_console_log"`
	DebugConsoleLogJSON bool `ini:"debug_console_log_json"`

	TypeMessage bool `ini:"type_message"`

	ListenIPv4Address string `ini:"listen_ipv4_address"`
	ListenPortNumber  int32  `ini:"listen_port_number"`

	// socket の待ち受け時間
	WebSocketReadTimeoutSec int32 `ini:"websocket_read_timeout_sec"`
	// pong が送られてこないためタイムアウトにするまでの時間
	WebSocketPongTimeoutSec int32 `ini:"websocket_pong_timeout_sec"`
	// ping 送信の時間間隔
	WebSocketPingIntervalSec int32 `ini:"websocket_ping_interval_sec"`

	// シグナリングからコピーする WebSocket の HTTP ヘッダー名
	CopyWebSocketHeaderNames []string `ini:"copy_websocket_header_names"`

	AuthnWebhookURL      string `ini:"authn_webhook_url"`
	DisconnectWebhookURL string `ini:"disconnect_webhook_url"`

	WebhookLogName           string `ini:"webhook_log_name"`
	WebhookRequestTimeoutSec int32  `ini:"webhook_request_timeout_sec"`

	ListenPrometheusIPv4Address string `ini:"listen_prometheus_ipv4_address"`
	ListenPrometheusPortNumber  int32  `ini:"listen_prometheus_port_number"`
}

func NewConfig(configFilePath string) (*Config, error) {
	config := new(Config)

	iniConfig, err := ini.InsensitiveLoad(configFilePath)
	if err != nil {
		return nil, err
	}

	if err := iniConfig.StrictMapTo(config); err != nil {
		return nil, err
	}

	if config.AuthnWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.AuthnWebhookURL); err != nil {
			return nil, err
		}
	}

	if config.DisconnectWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.DisconnectWebhookURL); err != nil {
			return nil, err
		}
	}

	setDefaultsConfig(config)

	return config, nil
}

func setDefaultsConfig(config *Config) {
	if config.LogDir == "" {
		config.LogDir = defaultLogDir
	}

	if config.LogName == "" {
		config.LogName = defaultLogName
	}

	if config.LogRotateMaxSize == 0 {
		config.LogRotateMaxSize = defaultLogRotateMaxSize
	}

	if config.LogRotateMaxBackups == 0 {
		config.LogRotateMaxBackups = defaultLogRotateMaxBackups
	}

	if config.LogRotateMaxAge == 0 {
		config.LogRotateMaxAge = defaultLogRotateMaxAge
	}

	if config.SignalingLogName == "" {
		config.SignalingLogName = defaultSignalingLogName
	}

	if config.SignalingLogFilters == nil {
		config.SignalingLogFilters = defaultSignalingLogFilters
	}

	if config.ListenIPv4Address == "" {
		config.ListenIPv4Address = defaultListenIPv4Address
	}

	if config.ListenPortNumber == 0 {
		config.ListenPortNumber = defaultListenPortNumber
	}

	if config.WebSocketReadTimeoutSec == 0 {
		config.WebSocketReadTimeoutSec = defaultWebSocketReadTimeoutSec
	}

	if config.WebSocketPongTimeoutSec == 0 {
		config.WebSocketPongTimeoutSec = defaultWebSocketPongTimeoutSec
	}

	if config.WebSocketPingIntervalSec == 0 {
		config.WebSocketPingIntervalSec = defaultWebSocketPingIntervalSec
	}

	if config.WebhookRequestTimeoutSec == 0 {
		config.WebhookRequestTimeoutSec = defaultWebhookRequestTimeout
	}

	if config.WebhookLogName == "" {
		config.WebhookLogName = defaultWebhookLogName
	}

	if config.ListenPrometheusIPv4Address == "" {
		config.ListenPrometheusIPv4Address = defaultListenPrometheusIPv4Address
	}

	if config.ListenPrometheusPortNumber == 0 {
		config.ListenPrometheusPortNumber = defaultListenPrometheusPortNumber
	}

}

func (c *Config) PrintConfig() {
	slog.Info("AyameConf", "debug", c.Debug)

	slog.Info("AyameConf", "log_dir", c.LogDir)
	slog.Info("AyameConf", "log_name", c.LogName)
	slog.Info("AyameConf", "log_stdout", c.LogStdout)

	slog.Info("AyameConf", "log_rotate_max_size", c.LogRotateMaxSize)
	slog.Info("AyameConf", "log_rotate_max_backups", c.LogRotateMaxBackups)
	slog.Info("AyameConf", "log_rotate_max_age", c.LogRotateMaxAge)
	slog.Info("AyameConf", "log_rotate_compress", c.LogRotateCompress)

	slog.Info("AyameConf", "signaling_log_name", c.SignalingLogName)
	slog.Info("AyameConf", "signaling_log_filters", c.SignalingLogFilters)

	slog.Info("AyameConf", "debug_console_log", c.DebugConsoleLog)
	slog.Info("AyameConf", "debug_console_log_json", c.DebugConsoleLogJSON)

	slog.Info("AyameConf", "listen_ipv4_address", c.ListenIPv4Address)
	slog.Info("AyameConf", "listen_port_number", c.ListenPortNumber)

	slog.Info("AyameConf", "copy_websocket_header_names", c.CopyWebSocketHeaderNames)

	slog.Info("AyameConf", "authn_webhook_url", c.AuthnWebhookURL)
	slog.Info("AyameConf", "disconnect_webhook_url", c.DisconnectWebhookURL)

	slog.Info("AyameConf", "webhook_log_name", c.WebhookLogName)
	slog.Info("AyameConf", "webhook_request_timeout_sec", c.WebhookRequestTimeoutSec)

	slog.Info("AyameConf", "prometheus_ipv4_address", c.ListenPrometheusIPv4Address)
	slog.Info("AyameConf", "prometheus_port", c.ListenPrometheusPortNumber)
}
