package ayame

import (
	"bufio"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	KB = 1 << 10
	MB = 1 << 20

	MetricsKey = "webhook_metrics"

	metricsSubsystem = "ayame"
)

// sizeBuckets は公開中のサイズ系ヒストグラムと同じバケット境界。
var sizeBuckets = []float64{
	1.0 * KB,
	2.0 * KB,
	5.0 * KB,
	10.0 * KB,
	100 * KB,
	500 * KB,
	1.0 * MB,
	2.5 * MB,
	5.0 * MB,
	10.0 * MB,
}

var httpLabels = []string{"code", "method", "host", "url"}

type Metrics struct {
	registry    *prometheus.Registry
	httpHandler http.Handler

	reqCnt *prometheus.CounterVec
	reqDur *prometheus.HistogramVec
	reqSz  *prometheus.HistogramVec
	resSz  *prometheus.HistogramVec

	webhookReqCnt   *prometheus.CounterVec
	webhookReqDur   *prometheus.HistogramVec
	webhookReqSz    *prometheus.HistogramVec
	webhookResSz    *prometheus.HistogramVec
	authnWebhookCnt *prometheus.CounterVec
}

func NewMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	// /metrics は従来どおり Go ランタイムとプロセスのメトリクスも返す。
	registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	m := &Metrics{registry: registry}
	m.reqCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: metricsSubsystem,
		Name:      "requests_total",
		Help:      "How many HTTP requests processed, partitioned by status code and HTTP method.",
	}, httpLabels)
	m.reqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricsSubsystem,
		Name:      "request_duration_seconds",
		Help:      "The HTTP request latencies in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, httpLabels)
	m.reqSz = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricsSubsystem,
		Name:      "request_size_bytes",
		Help:      "The HTTP request sizes in bytes.",
		Buckets:   sizeBuckets,
	}, httpLabels)
	m.resSz = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricsSubsystem,
		Name:      "response_size_bytes",
		Help:      "The HTTP response sizes in bytes.",
		Buckets:   sizeBuckets,
	}, httpLabels)
	m.webhookReqCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: metricsSubsystem,
		Name:      "webhook_requests_total",
		Help:      "How many HTTP requests.",
	}, httpLabels)
	m.webhookReqDur = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricsSubsystem,
		Name:      "webhook_request_duration_seconds",
		Help:      "The HTTP request latencies in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, httpLabels)
	m.webhookReqSz = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricsSubsystem,
		Name:      "webhook_request_message_size_bytes",
		Help:      "The HTTP request message sizes in bytes.",
		Buckets:   sizeBuckets,
	}, httpLabels)
	m.webhookResSz = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Subsystem: metricsSubsystem,
		Name:      "webhook_response_message_size_bytes",
		Help:      "The HTTP response message sizes in bytes.",
		Buckets:   sizeBuckets,
	}, httpLabels)
	m.authnWebhookCnt = prometheus.NewCounterVec(prometheus.CounterOpts{
		Subsystem: metricsSubsystem,
		Name:      "authn_webhook_responses_total",
		Help:      "How many AuthnWebhook responses.",
	}, []string{"code", "method", "host", "url", "allowed", "reason"})

	registry.MustRegister(
		m.reqCnt,
		m.reqDur,
		m.reqSz,
		m.resSz,
		m.webhookReqCnt,
		m.webhookReqDur,
		m.webhookReqSz,
		m.webhookResSz,
		m.authnWebhookCnt,
	)
	m.httpHandler = promhttp.InstrumentMetricHandler(
		registry,
		promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)
	return m
}

// middleware はシグナリング HTTP サーバのリクエストを計測する。
// /metrics は別サーバなのでここには来ない。
func (m *Metrics) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		reqSz := approximateRequestSize(r)
		rec := newStatusRecorder(w)
		next.ServeHTTP(rec, r)

		path := ""
		if r.URL != nil {
			path = r.URL.Path
		}
		m.observeHTTP(strconv.Itoa(rec.status), r.Method, r.Host, path, time.Since(start).Seconds(), reqSz, rec.size)
	})
}

func (m *Metrics) metricsHandler() http.Handler {
	return m.httpHandler
}

func (m *Metrics) observeHTTP(code, method, host, url string, elapsed float64, reqSz int, resSz int64) {
	m.reqCnt.WithLabelValues(code, method, host, url).Inc()
	m.reqDur.WithLabelValues(code, method, host, url).Observe(elapsed)
	m.reqSz.WithLabelValues(code, method, host, url).Observe(float64(reqSz))
	m.resSz.WithLabelValues(code, method, host, url).Observe(float64(resSz))
}

func (m *Metrics) IncWebhookReqCnt(code, method, host, url string) {
	m.webhookReqCnt.WithLabelValues(code, method, host, url).Inc()
}

func (m *Metrics) ObserveWebhookReqDur(code, method, host, url string, elapsed float64) {
	m.webhookReqDur.WithLabelValues(code, method, host, url).Observe(elapsed)
}

func (m *Metrics) ObserveWebhookResSz(code, method, host, url string, sz int64) {
	m.webhookResSz.WithLabelValues(code, method, host, url).Observe(float64(sz))
}

func (m *Metrics) ObserveWebhookReqSz(code, method, host, url string, sz int64) {
	m.webhookReqSz.WithLabelValues(code, method, host, url).Observe(float64(sz))
}

func (m *Metrics) IncAuthnWebhookCnt(code, method, host, url string, allowed bool, reason string) {
	m.authnWebhookCnt.WithLabelValues(code, method, host, url, strconv.FormatBool(allowed), reason).Inc()
}

// approximateRequestSize は ayame_request_size_bytes の値を維持するための近似。
func approximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s
}

// statusRecorder は応答ステータスとボディサイズを計測する。
// ステータス未設定は 200 とする。WebSocket の Upgrade は Hijack してから 101 を書くため、計測上の code は 200 になる。
type statusRecorder struct {
	http.ResponseWriter
	status      int
	size        int64
	wroteHeader bool
}

func newStatusRecorder(w http.ResponseWriter) *statusRecorder {
	return &statusRecorder{
		ResponseWriter: w,
		status:         http.StatusOK,
	}
}

func (w *statusRecorder) WriteHeader(code int) {
	if w.wroteHeader {
		return
	}
	w.status = code
	w.wroteHeader = true
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusRecorder) Write(b []byte) (int, error) {
	if !w.wroteHeader {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(b)
	w.size += int64(n)
	return n, err
}

func (w *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return h.Hijack()
}

func (w *statusRecorder) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}
