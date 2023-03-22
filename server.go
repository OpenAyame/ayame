package ayame

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Server struct {
	config *Config

	signalingLogger *zerolog.Logger
	webhookLogger   *zerolog.Logger

	http.Server
}

func NewServer(config *Config) (*Server, error) {
	signalingLogger, err := InitSignalingLogger(config)
	if err != nil {
		return nil, err
	}

	webhookLogger, err := InitWebhookLogger(config)
	if err != nil {
		return nil, err
	}

	s := &Server{
		config:          config,
		signalingLogger: signalingLogger,
		webhookLogger:   webhookLogger,
	}
	return s, nil
}

const readHeaderTimeout = 10 * time.Second

// TODO: echo 化したい
func (s *Server) Start(ctx context.Context) error {
	// URL の生成
	url := fmt.Sprintf("%s:%d", s.config.ListenIPv4Address, s.config.ListenPortNumber)

	// websocket server
	http.HandleFunc("/signaling", s.signalingHandler)
	http.HandleFunc("/.ok", s.okHandler)
	server := &http.Server{Addr: url, Handler: nil, ReadHeaderTimeout: readHeaderTimeout}

	ch := make(chan error)
	go func() {
		defer close(ch)
		if err := server.ListenAndServe(); err != nil {
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
