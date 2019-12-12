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

func authWebhookRequest(roomID string, clientID string, metadata interface{}, signalingKey *string) (*webhookResponse, error) {
	webhookReq := &webhookRequest{
		SignalingKey:  signalingKey,
		RoomID:        roomID,
		ClientID:      clientID,
		AuthnMetadata: metadata,
	}
	// TOOD(nakai): Timeout の引数を渡せるようにする
	respBody, err := postRequest(options.AuthWebhookURL, webhookReq)
	if err != nil {
		// TODO(nakai): ウェブフック失敗時に何故失敗したのかのログを追加する
		// TODO(nakai): ステータスコードなどもログとして出力するようにする
		return nil, err
	}
	webhookResp := webhookResponse{}
	err = json.Unmarshal(respBody, &webhookResp)
	if err != nil {
		// TODO(nakai): JSON がエラーになったのをログに追加する
		// roomID と clientID などを出力すること
		return nil, err
	}
	if !webhookResp.Allowed {
		logger.Info("authn webhook not allowed, resp=", &webhookResp)
		return &webhookResp, errors.New("Not Allowed")
	}
	logger.Info("auth webhook allowed, resp=", webhookResp)
	return &webhookResp, nil
}
