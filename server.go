package ayame

import (
	"context"
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

	mux := http.NewServeMux()
	mux.HandleFunc("GET /signaling", s.signalingHandler)
	mux.HandleFunc("GET /.ok", s.okHandler)

	s.Server = http.Server{
		Addr:              net.JoinHostPort(config.ListenIPv4Address, strconv.Itoa(int(config.ListenPortNumber))),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           metrics.middleware(mux),
	}

	metricsMux := http.NewServeMux()
	metricsMux.Handle("GET /metrics", metrics.metricsHandler())
	s.metricsServer = &http.Server{
		Addr:              net.JoinHostPort(config.ListenPrometheusIPv4Address, strconv.Itoa(int(config.ListenPrometheusPortNumber))),
		ReadHeaderTimeout: readHeaderTimeout,
		Handler:           metricsMux,
	}

	return s, nil
}

const readHeaderTimeout = 10 * time.Second

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
