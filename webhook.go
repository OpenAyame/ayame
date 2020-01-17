package main

import (
	"bytes"
	"encoding/json"
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
func (c *client) postRequest(u string, body interface{}) (*http.Response, error) {
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

	timeout := time.Duration(config.WebhookRequestTimeoutSec) * time.Second

	client := &http.Client{Timeout: timeout}
	return client.Do(req)
}

func (c *client) webhookLog(n string, v interface{}) {
	webhookLogger.Log().
		Str("roomId", c.roomID).
		Str("clientId", c.ID).
		Interface(n, v).
		Msg("")
}
