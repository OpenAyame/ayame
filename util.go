package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

// JSON HTTP Request をするだけのラッパー
func PostRequest(reqURL string, reqBody interface{}) ([]byte, error) {
	_, err := url.ParseRequestURI(reqURL)
	if err != nil {
		return nil, err
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		reqURL,
		bytes.NewBuffer([]byte(reqJSON)),
	)
	if err != nil {
		logger.Info("post request error, error=", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return respBody, nil
}

func trimOriginToHost(origin string) (*string, error) {
	url, err := url.Parse(origin)
	if err != nil {
		logger.Warning("origin parse error, origin=", origin)
		return nil, err
	}
	return &url.Host, nil
}
