package main

import (
	"encoding/json"
	"errors"
)

// webhook リクエスト
type WebhookRequest struct {
	Key    string `json:"key"`
	RoomId string `json:"room_id"`
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed    *bool  `json:"allowed"`
	WebhookUrl string `json:"auth_webhook_url"`
	Reason     string `json:"reason"`
}

// TODO(kdxu): 送信するデータを吟味する
type TwoAuthnRequest struct {
	Host          *string     `json:"host"`
	AuthnMetadata interface{} `json:"authn_metadata"`
}

type TwoAuthnResponse struct {
	Allowed       *bool       `json:"allowed"`
	AuthzMetadata interface{} `json:"authz_metadata"`
}

func AuthWebhookRequest(key string, roomId string, metadata interface{}, host string) (interface{}, error) {
	webhookReq := &WebhookRequest{Key: key, RoomId: roomId}
	respBytes, err := PostRequest(Options.AuthWebhookUrl, webhookReq)
	whResp := WebhookResponse{}
	err = json.Unmarshal(respBytes, &whResp)
	if err != nil {
		return nil, err
	}
	if !*whResp.Allowed {
		logger.Info("authn webhook not allowed, resp=", &whResp)
		return nil, errors.New("Not Allowed")
	}
	if whResp.WebhookUrl != "" {
		respBytes, err := PostRequest(whResp.WebhookUrl, &TwoAuthnRequest{Host: &host, AuthnMetadata: metadata})
		twoAuthnResp := TwoAuthnResponse{}
		err = json.Unmarshal(respBytes, &twoAuthnResp)
		if err != nil {
			return nil, err
		}
		if !*twoAuthnResp.Allowed {
			logger.Info("two authn webhook not allowed, resp=", &twoAuthnResp)
			return nil, errors.New("Not Allowed")
		}
		return twoAuthnResp.AuthzMetadata, nil
	}
	logger.Info("auth webhook allowed, resp=", whResp)
	return nil, nil
}
