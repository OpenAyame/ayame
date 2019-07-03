package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// JSON HTTP Request をするだけのラッパー
func PostRequest(url string, reqBody interface{}) ([]byte, error) {
	reqJson, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		url,
		bytes.NewBuffer([]byte(reqJson)),
	)
	if err != nil {
		logger.Info("post request error, error=", err)
		return nil, err
	}
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
