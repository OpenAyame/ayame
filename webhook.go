package main

import (
	"encoding/json"
	"errors"
)

type webhookRequest struct {
	SignalingKey  *string     `json:"signalingKey,omitempty"`
	RoomID        string      `json:"roomId"`
	ClientID      string      `json:"clientId"`
	AuthnMetadata interface{} `json:"authnMetadata"`
}

type webhookResponse struct {
	Allowed    bool          `json:"allowed"`
	IceServers []interface{} `json:"iceServers,omitempty"`
	Reason     string        `json:"reason"`
}

func authWebhookRequest(signalingKey *string, roomID string, clientID string, metadata interface{}) (*webhookResponse, error) {
	webhookReq := &webhookRequest{
		SignalingKey:  signalingKey,
		RoomID:        roomID,
		ClientID:      clientID,
		AuthnMetadata: metadata,
	}
	respBytes, err := PostRequest(options.AuthWebhookURL, webhookReq)
	if err != nil {
		return nil, err
	}
	webhookResp := webhookResponse{}
	err = json.Unmarshal(respBytes, &webhookResp)
	if err != nil {
		return nil, err
	}
	if !webhookResp.Allowed {
		logger.Info("authn webhook not allowed, resp=", &webhookResp)
		return &webhookResp, errors.New("Not Allowed")
	}
	logger.Info("auth webhook allowed, resp=", webhookResp)
	return &webhookResp, nil
}
