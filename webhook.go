package main

import (
	"encoding/json"
	"errors"
)

// webhook リクエスト
type WebhookRequest struct {
	SignalingKey  *string     `json:"signalingKey,omitempty"`
	RoomID        string      `json:"roomId"`
	ClientID      string      `json:"clientId"`
	AuthnMetadata interface{} `json:"authnMetadata"`
}

// webhook レスポンス
type WebhookResponse struct {
	Allowed    bool          `json:"allowed"`
	IceServers []interface{} `json:"iceServers,omitempty"`
	WebhookURL *string       `json:"authWebhookUrl,omitempty"`
	Reason     string        `json:"reason"`
}

func AuthWebhookRequest(signalingKey *string, roomID string, clientID string, metadata interface{}) (*WebhookResponse, error) {
	webhookReq := &WebhookRequest{SignalingKey: signalingKey, RoomID: roomID, ClientID: clientID, AuthnMetadata: metadata}
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
