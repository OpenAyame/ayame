package main

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	logrus "github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"time"
)

func NewLogger() *logrus.Logger {
	return &logrus.Logger{
		Formatter: &logrus.TextFormatter{},
	}
}

// logrus.logger に
// - ログローテ
// - ログディレクトリ
// - ログファイル名
// を設定する初期処理
// ayame 起動時に呼ばれる
// ログディレクトリおよびファイル名は起動時のオプションにて指定している
func setupLogger(l *logrus.Logger) error {
	if f, err := os.Stat(options.LogDir); os.IsNotExist(err) || !f.IsDir() {
		return err
	}
	level, err := logrus.ParseLevel(options.LogLevel)
	if err != nil {
		return err
	}
	l.SetLevel(level)

	logPath := fmt.Sprintf("%s/%s", options.LogDir, options.LogName)
	path, err := filepath.Abs(logPath + ".%Y%m%d")
	if err != nil {
		return err
	}
	rl, err := rotatelogs.New(path,
		rotatelogs.WithLinkName(logPath),
		rotatelogs.WithRotationTime(3600*time.Second),
	)
	if err != nil {
		return err
	}
	writer := io.MultiWriter(os.Stdout, rl)
	l.SetOutput(writer)

	l.Info("Setup log finished.")

	return nil
}
