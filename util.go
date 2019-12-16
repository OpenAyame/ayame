package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// JSON HTTP Request をするだけのラッパー
func postRequest(u string, body interface{}, timeout time.Duration) (*http.Response, error) {
	if _, err := url.ParseRequestURI(u); err != nil {
		return nil, err
	}
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
		logger.Info("post request error, error=", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: timeout}
	return client.Do(req)
}

func trimOriginToHost(origin string) (string, error) {
	url, err := url.Parse(origin)
	if err != nil {
		logger.Warning("origin parse error, origin=", origin)
		return "", err
	}
	return url.Host, nil
}
