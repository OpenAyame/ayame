package ayame

import (
	"fmt"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/shiguredo/lumberjack/v3"
)

func InitSignalingLogger(config *Config) (*zerolog.Logger, error) {
	if f, err := os.Stat(config.LogDir); os.IsNotExist(err) || !f.IsDir() {
		return nil, err
	}

	logPath := fmt.Sprintf("%s/%s", config.LogDir, config.SignalingLogName)

	// https://github.com/rs/zerolog/issues/77
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    logRotateMaxSize,
		MaxBackups: logRotateMaxBackups,
		MaxAge:     logRotateMaxAge,
		Compress:   true,
	}

	logger := zerolog.New(writer).With().Timestamp().Logger()

	return &logger, nil
}
