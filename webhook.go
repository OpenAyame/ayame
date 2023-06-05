package ayame

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"
)

type httpResponse struct {
	Status string      `json:"status"`
	Proto  string      `json:"proto"`
	Header http.Header `json:"header"`
	Body   string      `json:"body"`
}

// JSON HTTP Request をするだけのラッパー
func (c *connection) postRequest(u string, body interface{}) (*http.Response, error) {
	reqJSON, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		"POST",
		u,
		bytes.NewBuffer([]byte(reqJSON)),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	timeout := time.Duration(c.config.WebhookRequestTimeoutSec) * time.Second

	transport, err := c.transport()
	if err != nil {
		return nil, err
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   timeout,
	}
	return client.Do(req)
}

func (c *connection) transport() (*http.Transport, error) {
	// そもそも証明書のどちらかが設定されてなかったら、mTLS を利用しない
	if c.config.WebhookTLSFullchainFile == "" ||
		c.config.WebhookTLSPrivateKeyFile == "" ||
		c.config.WebhookTLSVerifyCacertFile == "" {
		return &http.Transport{}, nil
	}

	// まずは、証明書と秘密鍵を読み込む
	cert, err := tls.LoadX509KeyPair(c.config.WebhookTLSFullchainFile, c.config.WebhookTLSPrivateKeyFile)
	if err != nil {
		return nil, err
	}

	// 証明書プール（CA）を設定
	caCertPool := x509.NewCertPool()
	caCert, err := ioutil.ReadFile(c.config.WebhookTLSVerifyCacertFile) // CA証明書のパスを設定する
	if err != nil {
		return nil, err
	}
	caCertPool.AppendCertsFromPEM(caCert)

	// mTLSをサポートするためのTLS設定
	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		RootCAs:            caCertPool,
		InsecureSkipVerify: false,            // サーバ証明書を検証する場合はfalseに設定
		MinVersion:         tls.VersionTLS13, // 使用する最小TLSバージョンを指定する
	}

	// TLS設定をTransportに設定
	return &http.Transport{
		TLSClientConfig: tlsConfig,
	}, nil

}

func (c *connection) webhookLog(n string, v interface{}) {
	c.webhookLogger.Log().
		Str("roomId", c.roomID).
		Str("clientId", c.ID).
		Interface(n, v).
		Send()
}
