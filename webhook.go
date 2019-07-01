package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// webhook リクエスト
type WebhookRequest struct {
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed bool `json:"allowed"`
}

func authWebhookRequest() (interface{}, error) {
	webhookReq := &WebhookRequest{}
	reqJson, err := json.Marshal(webhookReq)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(
		"POST",
		Options.AuthWebhookUrl,
		bytes.NewBuffer([]byte(reqJson)),
	)
	if err != nil {
		logger.Info("auth webhook create request error, error=", err)
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Info("auth webhook post error, error=", err)
		return nil, err
	}
	whBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	whResp := WebhookResponse{}
	err = json.Unmarshal(whBody, &whResp)
	if err != nil {
		return nil, err
	}
	if !whResp.Allowed {
		logger.Info("auth webhook not allowed, resp=", whResp)
		return whResp, errors.New("Not Allowed")
	}
	logger.Info("auth webhook allowed, resp=", whResp)
	return whResp, nil
}
