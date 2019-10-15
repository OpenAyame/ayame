package main

import (
	"encoding/json"
	"errors"
)

// webhook リクエスト
type WebhookRequest struct {
	Key    *string `json:"key,omitempty"`
	RoomID string  `json:"room_id"`
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed       bool          `json:"allowed"`
	IceServers    []interface{} `json:"iceServers,omitempty"`
	WebhookURL    *string       `json:"auth_webhook_url,omitempty"`
	Reason        string        `json:"reason"`
	AuthzMetadata interface{}   `json:"authz_metadata"`
}

// TODO(kdxu): 送信するデータを吟味する
type TwoAuthnRequest struct {
	Host          *string     `json:"host"`
	AuthnMetadata interface{} `json:"authn_metadata"`
}

func AuthWebhookRequest(key *string, roomId string, metadata interface{}, host string) (*WebhookResponse, error) {
	webhookReq := &WebhookRequest{Key: key, RoomID: roomId}
	respBytes, err := PostRequest(Options.AuthWebhookURL, webhookReq)
	whResp := WebhookResponse{}
	err = json.Unmarshal(respBytes, &whResp)
	if err != nil {
		return nil, err
	}
	if !whResp.Allowed {
		logger.Info("authn webhook not allowed, resp=", &whResp)
		return &whResp, errors.New("Not Allowed")
	}
	if whResp.WebhookURL != nil {
		respBytes, err := PostRequest(*whResp.WebhookURL, &TwoAuthnRequest{Host: &host, AuthnMetadata: metadata})
		twoAuthnResp := WebhookResponse{IceServers: whResp.IceServers}
		err = json.Unmarshal(respBytes, &twoAuthnResp)
		if err != nil {
			return &whResp, err
		}
		if !twoAuthnResp.Allowed {
			logger.Info("two authn webhook not allowed, resp=", &twoAuthnResp)
			return &twoAuthnResp, errors.New("Not Allowed")
		}
		logger.Info("two authn webhook allowed, resp=", &twoAuthnResp)
		return &twoAuthnResp, nil
	}
	logger.Info("auth webhook allowed, resp=", whResp)
	return &whResp, nil
}
