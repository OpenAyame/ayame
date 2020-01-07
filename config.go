package main

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

	WebhookLogName        string `yaml:"webhook_log_name"`
	WebhookRequestTimeout int    `yaml:"webhook_request_timeout"`

	AllowOrigin string `yaml:"allow_origin"`
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

	if config.WebhookRequestTimeout == 0 {
		config.WebhookRequestTimeout = defaultWebhookRequestTimeout
	}
}
