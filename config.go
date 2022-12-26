package main

import (
	zlog "github.com/rs/zerolog/log"
)

const (
	defaultLogDir                = "."
	defaultLogName               = "ayame.log"
	defaultSignalingLogName      = "signaling.log"
	defaultWebhookLogName        = "webhook.log"
	defaultWebhookRequestTimeout = 5
)

type ayameConfig struct {
	Debug            bool   `yaml:"debug"`
	LogDir           string `yaml:"log_dir"`
	LogName          string `yaml:"log_name"`
	LogLevel         string `yaml:"log_level"`
	SignalingLogName string `yaml:"signaling_log_name"`

	ListenIPv4Address string `yaml:"listen_ipv4_address"`
	ListenPortNumber  int    `yaml:"listen_port_number"`

	AuthnWebhookURL      string `yaml:"authn_webhook_url"`
	DisconnectWebhookURL string `yaml:"disconnect_webhook_url"`

	WebhookLogName           string `yaml:"webhook_log_name"`
	WebhookRequestTimeoutSec uint   `yaml:"webhook_request_timeout_sec"`
}

func setDefaultsConfig() {
	if config.LogDir == "" {
		config.LogDir = defaultLogDir
	}

	if config.LogName == "" {
		config.LogDir = defaultLogName
	}

	if config.SignalingLogName == "" {
		config.SignalingLogName = defaultSignalingLogName
	}

	if config.WebhookLogName == "" {
		config.WebhookLogName = defaultWebhookLogName
	}

	if config.WebhookRequestTimeoutSec == 0 {
		config.WebhookRequestTimeoutSec = defaultWebhookRequestTimeout
	}

	zlog.Info().Bool("debug", config.Debug).Msg("AyameConf")
	zlog.Info().Str("log_dir", config.LogDir).Msg("AyameConf")
	zlog.Info().Str("log_name", config.LogName).Msg("AyameConf")
	zlog.Info().Str("log_level", config.LogLevel).Msg("AyameConf")
	zlog.Info().Str("signaling_log_name", config.SignalingLogName).Msg("AyameConf")
	zlog.Info().Str("listen_ipv4_address", config.ListenIPv4Address).Msg("AyameConf")
	zlog.Info().Int("listen_port_number", config.ListenPortNumber).Msg("AyameConf")
	zlog.Info().Str("authn_webhook_url", config.AuthnWebhookURL).Msg("AyameConf")
	zlog.Info().Str("disconnect_webhook_url", config.DisconnectWebhookURL).Msg("AyameConf")
	zlog.Info().Str("webhook_log_name", config.WebhookLogName).Msg("AyameConf")
	zlog.Info().Uint("webhook_request_timeout_sec", config.WebhookRequestTimeoutSec).Msg("AyameConf")
}
