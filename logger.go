package main

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

// logrus.logger に
// - ログローテ
// - ログディレクトリ
// - ログファイル名
// を設定する初期処理
// ayame 起動時に呼ばれる
// ログディレクトリおよびファイル名は起動時のオプションにて指定している
func setupLogger() *logrus.Logger {
	level, err := logrus.ParseLevel(options.LogLevel)
	if err != nil {
		log.Fatalf("Log level error %v", err)
	}
	logPath := fmt.Sprintf("%s/%s", options.LogDir, options.LogName)
	path, err := filepath.Abs(logPath + ".%Y%m%d")
	if err != nil {
		log.Fatalf("Log level error %v", err)
	}
	rl, err := rotatelogs.New(path,
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationTime(3600*time.Second),
	)
	if err != nil {
		log.Fatalf("Log level error %v", err)
	}
	out := io.MultiWriter(os.Stdout, rl)
	logger := logrus.Logger{
		Formatter: &logrus.TextFormatter{},
		Level:     level,
		Out:       out,
	}
	logger.Info("Setup log finished.")

	return &logger
}
