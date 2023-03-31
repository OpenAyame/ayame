package ayame

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	zlog "github.com/rs/zerolog/log"
)

type Server struct {
	config *Config

	signalingLogger *zerolog.Logger
	webhookLogger   *zerolog.Logger

	EchoPrometheus *echo.Echo
	Metrics        *Metrics

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

	e := echo.New()

	// URL の生成
	url := fmt.Sprintf("%s:%d", config.ListenIPv4Address, config.ListenPortNumber)

	s := &Server{
		config:          config,
		signalingLogger: signalingLogger,
		webhookLogger:   webhookLogger,
		Server: http.Server{
			Addr:              url,
			ReadHeaderTimeout: readHeaderTimeout,
			Handler:           e,
		},
	}

	// websocket server
	e.GET("/signaling", s.signalingHandler)
	e.GET("/.ok", s.okHandler)

	echoPrometheus := echo.New()
	echoPrometheus.HideBanner = true

	p := prometheus.NewPrometheus("ayame", nil, metricsList)
	e.Use(p.HandlerFunc)
	p.SetMetricsPath(echoPrometheus)

	s.EchoPrometheus = echoPrometheus

	m := NewMetrics()
	e.Use(m.AddMetricsMiddleware)

	s.Metrics = m

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
