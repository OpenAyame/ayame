package main

import (
	"encoding/json"
	"errors"
)

// webhook リクエスト
type WebhookRequest struct {
	Key      string `json:"key"`
	Metadata string `json:"auth_metadata"`
}

// TODO(kdxu): 送信するデータを吟味する
type TwoAuthnRequest struct {
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed    *bool  `json:"allowed"`
	WebhookUrl string `json:"webhook_url"`
	Reason     string `json:"reason"`
}

type TwoAuthnResponse struct {
	Allowed  *bool  `json:"allowed"`
	Metadata string `json:"authz_metadata"`
}

func authWebhookRequest(key string, metadata string) (interface{}, error) {
	webhookReq := &WebhookRequest{Key: key, Metadata: metadata}
	respBytes, err := PostRequest(Options.AuthWebhookUrl, webhookReq)
	whResp := WebhookResponse{}
	err = json.Unmarshal(respBytes, &whResp)
	if err != nil {
		return nil, err
	}
	if !*whResp.Allowed {
		logger.Info("auth webhook not allowed, resp=", &whResp)
		return whResp, errors.New("Not Allowed")
	}
	if whResp.WebhookUrl != "" {
		respBytes, err := PostRequest(whResp.WebhookUrl, &TwoAuthnRequest{})
		twoAuthnResp := TwoAuthnResponse{}
		err = json.Unmarshal(respBytes, &twoAuthnResp)
		if err != nil {
			return nil, err
		}
		if !*twoAuthnResp.Allowed {
			logger.Info("authz webhook not allowed, resp=", &twoAuthnResp)
			return whResp, errors.New("Not Allowed")
		}
		return twoAuthnResp, nil
	}
	logger.Info("auth webhook allowed, resp=", whResp)
	return whResp, nil
}
