package ayame

import (
	"time"

	"github.com/rs/zerolog"
)

const (
	ayameVersion = "2022.2.0"
	// timeout は暫定的に 10 sec
	readHeaderTimeout = 10 * time.Second
)

var (
	config          *Config
	signalingLogger *zerolog.Logger
	webhookLogger   *zerolog.Logger
)
