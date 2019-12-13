package main

import (
	"encoding/json"
)

type webhookRequest struct {
	SignalingKey  *string     `json:"signalingKey,omitempty"`
	RoomID        string      `json:"roomId"`
	ClientID      string      `json:"clientId"`
	AuthnMetadata interface{} `json:"authnMetadata"`
}

type webhookResponse struct {
	Allowed    *bool        `json:"allowed"`
	IceServers *[]iceServer `json:"iceServers,omitempty"`
	Reason     *string      `json:"reason",omitempty`
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
	if err := json.Unmarshal(respBody, &webhookResp); err != nil {
		// TODO(nakai): JSON がエラーになったのをログに追加する
		// roomID と clientID などを出力すること
		return nil, err
	}
	return &webhookResp, nil
}
