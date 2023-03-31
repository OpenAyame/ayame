package ayame

import (
	"github.com/labstack/echo-contrib/prometheus"
	"github.com/labstack/echo/v4"
	prom "github.com/prometheus/client_golang/prometheus"
)

const (
	KB = prometheus.KB
	MB = prometheus.MB
)

var (
	webhookReqDurBuckets = prom.DefBuckets
	webhookReqSzBuckets  = []float64{1.0 * KB, 2.0 * KB, 5.0 * KB, 10.0 * KB, 100 * KB, 500 * KB, 1.0 * MB, 2.5 * MB, 5.0 * MB, 10.0 * MB}
	webhookResSzBuckets  = []float64{1.0 * KB, 2.0 * KB, 5.0 * KB, 10.0 * KB, 100 * KB, 500 * KB, 1.0 * MB, 2.5 * MB, 5.0 * MB, 10.0 * MB}
)

var (
	authnWebhookReqCnt = &prometheus.Metric{
		ID:          "authnWebhookReqCnt",
		Name:        "authn_webhook_requests_total",
		Description: "How many HTTP requests.",
		Type:        "counter_vec",
		Args:        []string{"code", "method", "host", "url"}}
	authnWebhookReqDur = &prometheus.Metric{
		ID:          "authnWebhookReqDur",
		Name:        "authn_webhook_request_duration_seconds",
		Description: "The HTTP request latencies in seconds.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookReqDurBuckets}
	authnWebhookReqSz = &prometheus.Metric{
		ID:          "authnWebhookReqSz",
		Name:        "authn_webhook_request_size_bytes",
		Description: "The HTTP request sizes in bytes.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookReqSzBuckets}
	authnWebhookResSz = &prometheus.Metric{
		ID:          "authnWebhookResSz",
		Name:        "authn_webhook_response_size_bytes",
		Description: "The HTTP response sizes in bytes.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookResSzBuckets}

	disconnectWebhookReqCnt = &prometheus.Metric{
		ID:          "disconnectWebhookReqCnt",
		Name:        "disconnect_webhook_requests_total",
		Description: "How many HTTP requests.",
		Type:        "counter_vec",
		Args:        []string{"code", "method", "host", "url"}}
	disconnectWebhookReqDur = &prometheus.Metric{
		ID:          "disconnectWebhookReqDur",
		Name:        "disconnect_webhook_request_duration_seconds",
		Description: "The HTTP request latencies in seconds.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookReqDurBuckets}
	disconnectWebhookReqSz = &prometheus.Metric{
		ID:          "disconnectWebhookReqSz",
		Name:        "disconnect_webhook_request_size_bytes",
		Description: "The HTTP request sizes in bytes.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookReqSzBuckets}
	disconnectWebhookResSz = &prometheus.Metric{
		ID:          "disconnectWebhookResSz",
		Name:        "disconnect_webhook_response_size_bytes",
		Description: "The HTTP response sizes in bytes.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookResSzBuckets}

	metricsList = []*prometheus.Metric{
		authnWebhookReqCnt,
		authnWebhookReqDur,
		authnWebhookResSz,
		authnWebhookReqSz,

		disconnectWebhookReqCnt,
		disconnectWebhookReqDur,
		disconnectWebhookResSz,
		disconnectWebhookReqSz,
	}
)

type Metrics struct {
	AuthnWebhookReqCnt *prometheus.Metric
	AuthnWebhookReqDur *prometheus.Metric
	AuthnWebhookResSz  *prometheus.Metric
	AuthnWebhookReqSz  *prometheus.Metric

	DisconnectWebhookReqCnt *prometheus.Metric
	DisconnectWebhookReqDur *prometheus.Metric
	DisconnectWebhookResSz  *prometheus.Metric
	DisconnectWebhookReqSz  *prometheus.Metric
}

func NewMetrics() *Metrics {
	return &Metrics{
		AuthnWebhookReqCnt: authnWebhookReqCnt,
		AuthnWebhookReqDur: authnWebhookReqDur,
		AuthnWebhookResSz:  authnWebhookResSz,
		AuthnWebhookReqSz:  authnWebhookReqSz,

		DisconnectWebhookReqCnt: disconnectWebhookReqCnt,
		DisconnectWebhookReqDur: disconnectWebhookReqDur,
		DisconnectWebhookResSz:  disconnectWebhookResSz,
		DisconnectWebhookReqSz:  disconnectWebhookReqSz,
	}
}

var metricsKey = "webhook_metrics"

func (m *Metrics) AddMetricsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set(metricsKey, m)
		return next(c)
	}
}

func (m *Metrics) IncAuthnWebhookReqCnt(code, method, host, url string) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.AuthnWebhookReqCnt.MetricCollector.(*prom.CounterVec).With(labels).Inc()
}

func (m *Metrics) ObserveAuthnWebhookReqDur(code, method, host, url string, elapsed float64) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.AuthnWebhookReqDur.MetricCollector.(*prom.HistogramVec).With(labels).Observe(elapsed)
}

func (m *Metrics) ObserveAuthnWebhookResSz(code, method, host, url string, sz int) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.AuthnWebhookResSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) ObserveAuthnWebhookReqSz(code, method, host, url string, sz int) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.AuthnWebhookReqSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) IncDisconnectWebhookReqCnt(code, method, host, url string) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.DisconnectWebhookReqCnt.MetricCollector.(*prom.CounterVec).With(labels).Inc()
}

func (m *Metrics) ObserveDisconnectWebhookReqDur(code, method, host, url string, elapsed float64) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.DisconnectWebhookReqDur.MetricCollector.(*prom.HistogramVec).With(labels).Observe(elapsed)
}

func (m *Metrics) ObserveDisconnectWebhookResSz(code, method, host, url string, sz int) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.DisconnectWebhookResSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) ObserveDisconnectWebhookReqSz(code, method, host, url string, sz int) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.DisconnectWebhookReqSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}
