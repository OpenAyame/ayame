package main

import (
	"encoding/json"
	"io/ioutil"
)

// 後方互換性対応
// 次のリリースでは RoomId に揃える予定
// type authnWebhookRequest struct {
// 	RoomID        string       `json:"room_id"`
// 	ClientID      string       `json:"client_id"`
// 	SignalingKey  *string      `json:"signaling_key,omitempty"`
// 	Key           *string      `json:"key,omitempty"`
// 	AuthnMetadata *interface{} `json:"authn_metadata,omitempty"`
// 	AyameClient   *string      `json:"ayame_client,omitempty"`
// 	Libwebrtc     *string      `json:"libwebrtc,omitempty"`
// 	Environment   *string      `json:"environment,omitempty"`
// }

type authnWebhookRequest struct {
	RoomID        string       `json:"roomId"`
	ClientID      string       `json:"clientId"`
	SignalingKey  *string      `json:"signalingKey,omitempty"`
	Key           *string      `json:"key,omitempty"`
	AuthnMetadata *interface{} `json:"authnMetadata,omitempty"`
	AyameClient   *string      `json:"ayameClient,omitempty"`
	Libwebrtc     *string      `json:"libwebrtc,omitempty"`
	Environment   *string      `json:"environment,omitempty"`
}

type authnWebhookResponse struct {
	Allowed       *bool        `json:"allowed"`
	IceServers    *[]iceServer `json:"iceServers,omitempty"`
	Reason        *string      `json:"reason,omitempty"`
	AuthzMetadata *interface{} `json:"authzMetadata,omitempty"`
}

func (c *client) authnWebhook() (*authnWebhookResponse, error) {
	if config.AuthnWebhookURL == "" {
		var allowed = true
		authnWebhookResponse := &authnWebhookResponse{Allowed: &allowed}
		return authnWebhookResponse, nil
	}

	req := &authnWebhookRequest{
		RoomID:       c.roomID,
		ClientID:     c.ID,
		SignalingKey: c.signalingKey,
		// 後方互換性対応
		Key:           c.signalingKey,
		AuthnMetadata: c.authnMetadata,
	}

	resp, err := c.postRequest(config.AuthnWebhookURL, req)
	if err != nil {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Err(err).
			Caller().
			Msg("AuthnWebhookError")
		return nil, err
	}
	// http://ikawaha.hateblo.jp/entry/2015/06/07/074155
	defer resp.Body.Close()

	c.webhookLog("authnReq", req)
	c.webhookLog("authnResp", resp)

	// 200 以外で返ってきたときはエラーとする
	if resp.StatusCode != 200 {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Interface("resp", resp).
			Caller().
			Msg("AuthnWebhookResponseError")
		return nil, errAuthnWebhookResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Bytes("body", body).
			Err(err).
			Caller().
			Msg("AuthnWebhookResponseError")
		return nil, err
	}

	authnWebhookResponse := authnWebhookResponse{}
	if err := json.Unmarshal(body, &authnWebhookResponse); err != nil {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Err(err).
			Caller().
			Msg("AuthnWebhookResponseError")
		return nil, errAuthnWebhookResponse
	}

	return &authnWebhookResponse, nil
}
