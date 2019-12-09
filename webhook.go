package main

import (
	"encoding/json"
	"errors"
)

// webhook リクエスト
type WebhookRequest struct {
	Key    *string `json:"key,omitempty"`
	RoomID string  `json:"roomId"`
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed    bool          `json:"allowed"`
	IceServers []interface{} `json:"iceServers,omitempty"`
	WebhookURL *string       `json:"authWebhookUrl,omitempty"`
	Reason     string        `json:"reason"`
}

func AuthWebhookRequest(key *string, roomID string, metadata interface{}) (*WebhookResponse, error) {
	webhookReq := &WebhookRequest{Key: key, RoomID: roomID}
	respBytes, err := PostRequest(Options.AuthWebhookURL, webhookReq)
	if err != nil {
		return nil, err
	}
	whResp := WebhookResponse{}
	err = json.Unmarshal(respBytes, &whResp)
	if err != nil {
		return nil, err
	}
	if !whResp.Allowed {
		logger.Info("authn webhook not allowed, resp=", &whResp)
		return &whResp, errors.New("Not Allowed")
	}
	logger.Info("auth webhook allowed, resp=", whResp)
	return &whResp, nil
}
