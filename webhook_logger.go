package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

func initWebhookLogger() (*zerolog.Logger, error) {

	if f, err := os.Stat(config.LogDir); os.IsNotExist(err) || !f.IsDir() {
		return nil, err
	}

	logPath := fmt.Sprintf("%s/%s", config.LogDir, config.WebhookLogName)

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

	var writers io.Writer
	// デバッグが有効な時はコンソールにもだす
	if config.Debug {
		writers = io.MultiWriter(os.Stdout, writer)
	}

	logger := zerolog.New(writers).With().Timestamp().Logger()

	return &logger, nil
}
