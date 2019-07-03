package main

import (
	"encoding/json"
	"errors"
)

// webhook リクエスト
type WebhookRequest struct {
	Key      string `json:key`
	Metadata string `json:"authn_metadata"`
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed bool   `json:"allowed"`
	Reason  string `json:"reason"`
}

func authWebhookRequest(key string, metadata string) (interface{}, error) {
	webhookReq := &WebhookRequest{Key: key, Metadata: metadata}
	respBytes, err := PostRequest(Options.AuthWebhookUrl, webhookReq)
	whResp := WebhookResponse{}
	err = json.Unmarshal(respBytes, &whResp)
	if err != nil {
		return nil, err
	}
	if !whResp.Allowed {
		logger.Info("auth webhook not allowed, resp=", &whResp)
		return whResp, errors.New("Not Allowed")
	}
	logger.Info("auth webhook allowed, resp=", whResp)
	return whResp, nil
}
