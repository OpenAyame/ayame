package ayame

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T) *Server {
	t.Helper()

	cfg := &Config{
		LogDir:                      t.TempDir(),
		LogName:                     "ayame.jsonl",
		SignalingLogName:            "signaling.jsonl",
		WebhookLogName:              "webhook.jsonl",
		LogRotateMaxSize:            1,
		LogRotateMaxBackups:         1,
		LogRotateMaxAge:             1,
		ListenIPv4Address:           "127.0.0.1",
		ListenPortNumber:            3000,
		ListenPrometheusIPv4Address: "127.0.0.1",
		ListenPrometheusPortNumber:  4000,
		WebSocketReadTimeoutSec:     30,
		WebSocketPongTimeoutSec:     30,
		WebSocketPingIntervalSec:    30,
	}
	s, err := NewServer(cfg)
	require.NoError(t, err)
	return s
}

func scrapeMetrics(t *testing.T, s *Server) string {
	t.Helper()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	s.metricsServer.Handler.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	return rec.Body.String()
}

func TestServerRoutesAndMetrics(t *testing.T) {
	s := newTestServer(t)

	assert.Equal(t, "127.0.0.1:3000", s.Addr)
	assert.Equal(t, "127.0.0.1:4000", s.metricsServer.Addr)
	assert.Equal(t, readHeaderTimeout, s.ReadHeaderTimeout)
	assert.Equal(t, readHeaderTimeout, s.metricsServer.ReadHeaderTimeout)

	req := httptest.NewRequest(http.MethodGet, "/.ok", nil)
	req.Host = "ayame.local"
	rec := httptest.NewRecorder()
	s.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Body.String())

	req = httptest.NewRequest(http.MethodPost, "/.ok", nil)
	req.Host = "ayame.local"
	rec = httptest.NewRecorder()
	s.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/missing", nil)
	req.Host = "ayame.local"
	rec = httptest.NewRecorder()
	s.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/signaling", nil)
	req.Host = "ayame.local"
	rec = httptest.NewRecorder()
	s.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// /metrics はシグナリングサーバには出さない。
	req = httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rec = httptest.NewRecorder()
	s.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/.ok", nil)
	rec = httptest.NewRecorder()
	s.metricsServer.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)

	req = httptest.NewRequest(http.MethodPost, "/metrics", nil)
	rec = httptest.NewRecorder()
	s.metricsServer.Handler.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

	s.Metrics.IncWebhookReqCnt("204", "POST", "hook.example", "/disconnect")
	s.Metrics.ObserveWebhookReqDur("204", "POST", "hook.example", "/disconnect", 0.01)
	s.Metrics.ObserveWebhookReqSz("204", "POST", "hook.example", "/disconnect", 128)
	s.Metrics.ObserveWebhookResSz("204", "POST", "hook.example", "/disconnect", 64)
	s.Metrics.IncAuthnWebhookCnt("200", "POST", "hook.example", "/authn", false, "deny")

	body := scrapeMetrics(t, s)
	assert.Contains(t, body, "go_goroutines")
	assert.Contains(t, body, `ayame_requests_total{code="200",host="ayame.local",method="GET",url="/.ok"} 1`)
	assert.Contains(t, body, `ayame_requests_total{code="405",host="ayame.local",method="POST",url="/.ok"} 1`)
	assert.Contains(t, body, `ayame_requests_total{code="404",host="ayame.local",method="GET",url="/missing"} 1`)
	assert.Contains(t, body, `ayame_requests_total{code="400",host="ayame.local",method="GET",url="/signaling"} 1`)
	assert.Contains(t, body, `ayame_request_duration_seconds_count{code="200",host="ayame.local",method="GET",url="/.ok"} 1`)
	assert.Contains(t, body, `ayame_webhook_requests_total{code="204",host="hook.example",method="POST",url="/disconnect"} 1`)
	assert.Contains(t, body, `ayame_webhook_request_duration_seconds_count{code="204",host="hook.example",method="POST",url="/disconnect"} 1`)
	assert.Contains(t, body, `ayame_webhook_request_message_size_bytes_count{code="204",host="hook.example",method="POST",url="/disconnect"} 1`)
	assert.Contains(t, body, `ayame_webhook_response_message_size_bytes_count{code="204",host="hook.example",method="POST",url="/disconnect"} 1`)
	assert.Contains(t, body, `ayame_authn_webhook_responses_total{allowed="false",code="200",host="hook.example",method="POST",reason="deny",url="/authn"} 1`)
}

func TestServerWebSocketUpgrade(t *testing.T) {
	s := newTestServer(t)
	ts := httptest.NewServer(s.Handler)
	t.Cleanup(ts.Close)

	wsURL := "ws" + strings.TrimPrefix(ts.URL, "http") + "/signaling"
	dialer := websocket.Dialer{HandshakeTimeout: 2 * time.Second}
	conn, resp, err := dialer.Dial(wsURL, nil)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	require.NoError(t, err)
	t.Cleanup(func() {
		conn.Close()
	})
	assert.Equal(t, http.StatusSwitchingProtocols, resp.StatusCode)

	host := strings.TrimPrefix(ts.URL, "http://")
	want := `ayame_requests_total{code="200",host="` + host + `",method="GET",url="/signaling"} 1`
	deadline := time.Now().Add(2 * time.Second)
	var body string
	for {
		body = scrapeMetrics(t, s)
		if strings.Contains(body, want) {
			return
		}
		if time.Now().After(deadline) {
			t.Fatalf("metric %q not found\n%s", want, body)
		}
		time.Sleep(5 * time.Millisecond)
	}
}
