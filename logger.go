package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// megabytes
	logRotateMaxSize    = 10
	logRotateMaxBackups = 5
	//days
	logRotateMaxAge = 30
)

func initLogger() (*zerolog.Logger, error) {
	if f, err := os.Stat(config.LogDir); os.IsNotExist(err) || !f.IsDir() {
		return nil, err
	}

	logPath := fmt.Sprintf("%s/%s", config.LogDir, config.LogName)

	writer := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    logRotateMaxSize,
		MaxBackups: logRotateMaxBackups,
		MaxAge:     logRotateMaxAge,
		Compress:   true,
	}

	// https://github.com/rs/zerolog/issues/77
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	var writers io.Writer
	// デバッグが有効な時はコンソールにもだす
	if config.Debug {
		writers = io.MultiWriter(os.Stdout, writer)
	}

	output := zerolog.ConsoleWriter{Out: writers, TimeFormat: "2006-01-02 15:04:05.000000Z"}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%s]", i))
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}

	logger := zerolog.New(output).With().Timestamp().Logger()

	return &logger, nil
}
