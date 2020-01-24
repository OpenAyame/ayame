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

	var logLevel zerolog.Level
	// debug: true の場合の log_level は debug で固定
	if config.Debug {
		logLevel = zerolog.DebugLevel
	} else {
		switch config.LogLevel {
		case "info":
			logLevel = zerolog.InfoLevel
		case "error":
			logLevel = zerolog.ErrorLevel
		case "fatal":
			logLevel = zerolog.FatalLevel
		case "debug":
			logLevel = zerolog.DebugLevel
		default:
			return nil, errConfigInvalidLogLevel
		}
	}

	zerolog.SetGlobalLevel(logLevel)

	// https://github.com/rs/zerolog/issues/77
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	zerolog.TimeFieldFormat = time.RFC3339Nano

	var writers io.Writer
	output := zerolog.ConsoleWriter{Out: writer, NoColor: true, TimeFormat: "2006-01-02 15:04:05.000000Z"}
	format(&output)
	// デバッグが有効な時はコンソールにもだす
	if config.Debug {
		stdout := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05.000000Z"}
		format(&stdout)
		writers = io.MultiWriter(stdout, output)
	} else {
		writers = output
	}

	logger := zerolog.New(writers).With().Timestamp().Logger()

	return &logger, nil
}

func format(w *zerolog.ConsoleWriter) {
	w.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("[%s]", i))
	}
	w.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s=", i)
	}
	w.FormatFieldValue = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
}
