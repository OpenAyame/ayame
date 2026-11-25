package ayame

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Server struct {
	config *Config

	signalingLogger *zerolog.Logger
	webhookLogger   *zerolog.Logger

	Metrics *Metrics

	http.Server
	metricsServer *http.Server
}

func NewServer(config *Config) (*Server, error) {
	// signaling.jsonl ロガーを用意
	signalingLogger, err := NewLogger(config, config.SignalingLogName, "signaling")
	if err != nil {
		return nil, err
	}

	// webhook.jsonl ロガーを用意
	webhookLogger, err := NewLogger(config, config.WebhookLogName, "webhook")
	if err != nil {
		return nil, err
	}

	metrics := NewMetrics()

	s := &Server{
		config:          config,
		signalingLogger: signalingLogger,
		webhookLogger:   webhookLogger,
		Metrics:         metrics,
	}

	router := newMethodRouter(map[string]http.Handler{
		"/signaling": http.HandlerFunc(s.signalingHandler),
		"/.ok":       http.HandlerFunc(s.okHandler),
	})
	s.Server = http.Server{
		Addr:              net.JoinHostPort(config.ListenIPv4Address, strconv.Itoa(int(config.ListenPortNumber))),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           metrics.middleware(router),
	}

	// メトリクス用サーバは echo と同様、ReadHeaderTimeout を付けない。
	s.metricsServer = &http.Server{
		Addr: net.JoinHostPort(config.ListenPrometheusIPv4Address, strconv.Itoa(int(config.ListenPrometheusPortNumber))),
		Handler: newMethodRouter(map[string]http.Handler{
			"/metrics": metrics.metricsHandler(),
		}),
	}

	return s, nil
}

const (
	readHeaderTimeout = 10 * time.Second

	// echo が GET ルートに出す Allow。HEAD は含めない。
	allowOptionsGet = "OPTIONS, GET"
)

// methodRouter は echo と同じメソッド扱いをする。
// 登録パスの GET だけを処理し、OPTIONS は 204、それ以外は 405 を返す。
// 未登録パスは 404 を返す。
type methodRouter struct {
	gets map[string]http.Handler
}

func newMethodRouter(gets map[string]http.Handler) http.Handler {
	return methodRouter{gets: gets}
}

func (m methodRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, ok := m.gets[requestPath(r)]
	if !ok {
		writeEchoStatus(w, r, http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		handler.ServeHTTP(w, r)
	case http.MethodOptions:
		w.Header().Set("Allow", allowOptionsGet)
		w.WriteHeader(http.StatusNoContent)
	default:
		w.Header().Set("Allow", allowOptionsGet)
		writeEchoStatus(w, r, http.StatusMethodNotAllowed)
	}
}

// requestPath は echo の GetPath と同じく、RawPath があればそれを使う。
func requestPath(r *http.Request) string {
	if r.URL == nil {
		return ""
	}
	if r.URL.RawPath != "" {
		return r.URL.RawPath
	}
	return r.URL.Path
}

// writeEchoStatus は echo の DefaultHTTPErrorHandler と同じ JSON を返す。
// HEAD はボディを書かない。
func writeEchoStatus(w http.ResponseWriter, r *http.Request, code int) {
	if r.Method == http.MethodHead {
		w.WriteHeader(code)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"message": http.StatusText(code),
	})
}

func (s *Server) Start(ctx context.Context) error {
	ch := make(chan error)
	go func() {
		defer close(ch)
		if err := s.ListenAndServe(); err != nil {
			ch <- err
		}
	}()

	defer func() {
		if err := s.Shutdown(ctx); err != nil {
			zlog.Error().Err(err).Send()
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case err := <-ch:
		return err
	}
}

// StartMetrics は Prometheus の /metrics を配信する HTTP サーバを起動する。
func (s *Server) StartMetrics() error {
	return s.metricsServer.ListenAndServe()
}
