package ayame

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLogger(config *Config) error {
	// slog doesn't have global message/timestamp field name configuration
	// These would need to be handled in custom handlers if needed
	return nil
}

func NewLogger(config *Config, logFilename string, logDomain string) (*slog.Logger, error) {
	var handler slog.Handler
	
	// デバッグコンソールログを出力する
	// デバッグコンソールには Caller を出力する
	if config.Debug && config.DebugConsoleLog {
		// デバッグコンソールを JSON 形式で出力
		if config.DebugConsoleLogJSON {
			opts := &slog.HandlerOptions{
				Level: slog.LevelDebug,
				AddSource: true,
			}
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			// デバッグコンソールをヒューマンリーダブルな形式で出力
			opts := &slog.HandlerOptions{
				Level: slog.LevelDebug,
				AddSource: true,
				ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
					if a.Key == slog.SourceKey {
						source := a.Value.Any().(*slog.Source)
						return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", filepath.Base(source.File), source.Line))
					}
					return a
				},
			}
			handler = newPrettyHandler(os.Stdout, opts)
		}
	} else {
		// 標準出力にログを出力する
		level := slog.LevelInfo
		if config.Debug {
			level = slog.LevelDebug
		}
		
		if config.LogStdout {
			opts := &slog.HandlerOptions{
				Level: level,
			}
			handler = slog.NewJSONHandler(os.Stdout, opts)
		} else {
			// ログファイルを出力する
			if f, err := os.Stat(config.LogDir); os.IsNotExist(err) || !f.IsDir() {
				return nil, err
			}

			logPath := fmt.Sprintf("%s/%s", config.LogDir, logFilename)

			lumberjackLogger := &lumberjack.Logger{
				Filename:   logPath,
				MaxSize:    config.LogRotateMaxSize,
				MaxBackups: config.LogRotateMaxBackups,
				MaxAge:     config.LogRotateMaxAge,
				Compress:   config.LogRotateCompress,
			}
			
			opts := &slog.HandlerOptions{
				Level: level,
			}
			handler = slog.NewJSONHandler(lumberjackLogger, opts)
		}
	}

	logger := slog.New(handler).With("domain", logDomain)
	return logger, nil
}

// prettyHandler implements a custom handler for human-readable output
type prettyHandler struct {
	opts *slog.HandlerOptions
	w    io.Writer
}

func newPrettyHandler(w io.Writer, opts *slog.HandlerOptions) *prettyHandler {
	return &prettyHandler{
		opts: opts,
		w:    w,
	}
}

func (h *prettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return level >= h.opts.Level.Level()
}

func (h *prettyHandler) Handle(ctx context.Context, r slog.Record) error {
	const (
		reset     = "\x1b[0m"
		darkGray  = "\x1b[90m"
		red       = "\x1b[31m"
		green     = "\x1b[32m"
		yellow    = "\x1b[33m"
		blue      = "\x1b[34m"
		cyan      = "\x1b[36m"
		white     = "\x1b[37m"
	)

	// Format timestamp
	timestamp := r.Time.UTC().Format(time.RFC3339Nano)
	fmt.Fprintf(h.w, "%s%s%s ", darkGray, timestamp, reset)

	// Format level
	var levelColor, levelText string
	switch r.Level {
	case slog.LevelDebug:
		levelColor = blue
		levelText = "DEBUG"
	case slog.LevelInfo:
		levelColor = green
		levelText = "INFO"
	case slog.LevelWarn:
		levelColor = yellow
		levelText = "WARN"
	case slog.LevelError:
		levelColor = red
		levelText = "ERROR"
	default:
		levelColor = white
		levelText = r.Level.String()
	}
	fmt.Fprintf(h.w, "%s[%s]%s ", levelColor, levelText, reset)

	// Add source if enabled
	if h.opts.AddSource && r.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{r.PC})
		f, _ := fs.Next()
		if f.File != "" {
			fmt.Fprintf(h.w, "[%s:%d] ", filepath.Base(f.File), f.Line)
		}
	}

	// Format message
	if r.Message != "" {
		fmt.Fprintf(h.w, "%s | ", r.Message)
	}

	// Format attributes
	first := true
	r.Attrs(func(a slog.Attr) bool {
		if !first {
			fmt.Fprint(h.w, " ")
		}
		first = false
		fmt.Fprintf(h.w, "%s%s=%s%v", cyan, a.Key, reset, a.Value)
		return true
	})

	fmt.Fprintln(h.w)
	return nil
}

func (h *prettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// For simplicity, we're not implementing this fully
	// In a production implementation, you'd want to store and apply these attrs
	return h
}

func (h *prettyHandler) WithGroup(name string) slog.Handler {
	// For simplicity, we're not implementing this fully
	return h
}