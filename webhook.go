package main

import (
	"encoding/json"
	"io/ioutil"
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
	Reason     *string      `json:"reason,omitempty"`
}

func authWebhookRequest(roomID string, clientID string, authnMetadata interface{}, signalingKey *string) (*webhookResponse, error) {
	req := &webhookRequest{
		SignalingKey:  signalingKey,
		RoomID:        roomID,
		ClientID:      clientID,
		AuthnMetadata: authnMetadata,
	}
	// TOOD(nakai): Timeout の引数を渡せるようにする
	resp, err := postRequest(options.AuthWebhookURL, req)
	if err != nil {
		// TODO(nakai): ウェブフック失敗時に何故失敗したのかのログを追加する
		// TODO(nakai): ステータスコードなどもログとして出力するようにする
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	webhookResponse := webhookResponse{}
	if err := json.Unmarshal(body, &webhookResponse); err != nil {
		// TODO(nakai): JSON がエラーになったのをログに追加する
		// roomID と clientID などを出力すること
		return nil, err
	}
	return &webhookResponse, nil
}
