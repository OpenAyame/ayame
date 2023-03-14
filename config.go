package ayame

import (
	"net/url"

	zlog "github.com/rs/zerolog/log"
	"gopkg.in/ini.v1"
)

const (
	defaultLogDir                = "."
	defaultLogName               = "ayame.log"
	defaultSignalingLogName      = "signaling.log"
	defaultWebhookLogName        = "webhook.log"
	defaultWebhookRequestTimeout = 5
)

type Config struct {
	Debug            bool   `ini:"debug"`
	LogDir           string `ini:"log_dir"`
	LogName          string `ini:"log_name"`
	LogLevel         string `ini:"log_level"`
	SignalingLogName string `ini:"signaling_log_name"`

	ListenIPv4Address string `ini:"listen_ipv4_address"`
	ListenPortNumber  int32  `ini:"listen_port_number"`

	AuthnWebhookURL      string `ini:"authn_webhook_url"`
	DisconnectWebhookURL string `ini:"disconnect_webhook_url"`

	WebhookLogName           string `ini:"webhook_log_name"`
	WebhookRequestTimeoutSec int32  `ini:"webhook_request_timeout_sec"`
}

func InitConfig(iniConfig ini.File, config *Config) error {
	if err := iniConfig.MapTo(config); err != nil {
		return err
	}

	if config.AuthnWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.AuthnWebhookURL); err != nil {
			return err
		}
	}

	if config.DisconnectWebhookURL != "" {
		if _, err := url.ParseRequestURI(config.DisconnectWebhookURL); err != nil {
			return err
		}
	}

	return nil
}

func SetDefaultsConfig(config Config) {
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
	zlog.Info().Int32("listen_port_number", config.ListenPortNumber).Msg("AyameConf")
	zlog.Info().Str("authn_webhook_url", config.AuthnWebhookURL).Msg("AyameConf")
	zlog.Info().Str("disconnect_webhook_url", config.DisconnectWebhookURL).Msg("AyameConf")
	zlog.Info().Str("webhook_log_name", config.WebhookLogName).Msg("AyameConf")
	zlog.Info().Int32("webhook_request_timeout_sec", config.WebhookRequestTimeoutSec).Msg("AyameConf")
}
