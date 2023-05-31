package ayame

import (
	"fmt"

	"github.com/labstack/echo-contrib/prometheus"
	prom "github.com/prometheus/client_golang/prometheus"
)

const (
	KB = prometheus.KB
	MB = prometheus.MB

	MetricsKey = "webhook_metrics"
)

var (
	webhookReqDurBuckets = prom.DefBuckets
	webhookReqSzBuckets  = []float64{1.0 * KB, 2.0 * KB, 5.0 * KB, 10.0 * KB, 100 * KB, 500 * KB, 1.0 * MB, 2.5 * MB, 5.0 * MB, 10.0 * MB}
	webhookResSzBuckets  = []float64{1.0 * KB, 2.0 * KB, 5.0 * KB, 10.0 * KB, 100 * KB, 500 * KB, 1.0 * MB, 2.5 * MB, 5.0 * MB, 10.0 * MB}
)

var (
	webhookReqCnt = &prometheus.Metric{
		ID:          "webhookReqCnt",
		Name:        "webhook_requests_total",
		Description: "How many HTTP requests.",
		Type:        "counter_vec",
		Args:        []string{"code", "method", "host", "url"}}
	webhookReqDur = &prometheus.Metric{
		ID:          "webhookReqDur",
		Name:        "webhook_request_duration_seconds",
		Description: "The HTTP request latencies in seconds.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookReqDurBuckets}
	webhookReqSz = &prometheus.Metric{
		ID:          "webhookReqSz",
		Name:        "webhook_request_message_size_bytes",
		Description: "The HTTP request message sizes in bytes.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookReqSzBuckets}
	webhookResSz = &prometheus.Metric{
		ID:          "webhookResSz",
		Name:        "webhook_response_message_size_bytes",
		Description: "The HTTP response message sizes in bytes.",
		Args:        []string{"code", "method", "host", "url"},
		Type:        "histogram_vec",
		Buckets:     webhookResSzBuckets}
	authnWebhookCnt = &prometheus.Metric{
		ID:          "authnWebhookRespCnt",
		Name:        "authn_webhook_responses_total",
		Description: "How many AuthnWebhook responses.",
		Type:        "counter_vec",
		Args:        []string{"code", "method", "host", "url", "allowed", "reason"}}

	metricsList = []*prometheus.Metric{
		webhookReqCnt,
		webhookReqDur,
		webhookResSz,
		webhookReqSz,
		authnWebhookCnt,
	}
)

type Metrics struct {
	WebhookReqCnt   *prometheus.Metric
	WebhookReqDur   *prometheus.Metric
	WebhookResSz    *prometheus.Metric
	WebhookReqSz    *prometheus.Metric
	AuthnWebhookCnt *prometheus.Metric
}

func NewMetrics() *Metrics {
	return &Metrics{
		WebhookReqCnt:   webhookReqCnt,
		WebhookReqDur:   webhookReqDur,
		WebhookResSz:    webhookResSz,
		WebhookReqSz:    webhookReqSz,
		AuthnWebhookCnt: authnWebhookCnt,
	}
}

func (m *Metrics) IncWebhookReqCnt(code, method, host, url string) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.WebhookReqCnt.MetricCollector.(*prom.CounterVec).With(labels).Inc()
}

func (m *Metrics) ObserveWebhookReqDur(code, method, host, url string, elapsed float64) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.WebhookReqDur.MetricCollector.(*prom.HistogramVec).With(labels).Observe(elapsed)
}

func (m *Metrics) ObserveWebhookResSz(code, method, host, url string, sz int64) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.WebhookResSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) ObserveWebhookReqSz(code, method, host, url string, sz int64) {
	labels := prom.Labels{
		"code":   code,
		"method": method,
		"host":   host,
		"url":    url,
	}
	m.WebhookReqSz.MetricCollector.(*prom.HistogramVec).With(labels).Observe(float64(sz))
}

func (m *Metrics) IncAuthnWebhookCnt(code, method, host, url string, allowed bool, reason string) {
	labels := prom.Labels{
		"code":    code,
		"method":  method,
		"host":    host,
		"url":     url,
		"allowed": fmt.Sprintf("%v", allowed),
		"reason":  reason,
	}
	m.AuthnWebhookCnt.MetricCollector.(*prom.CounterVec).With(labels).Inc()
}
