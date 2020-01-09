package main

import (
	"io/ioutil"
)

type disconnectWebhookRequest struct {
	RoomID   string `json:"roomId"`
	ClientID string `json:"clientId"`
}

func (c *client) disconnectWebhook() error {
	if config.DisconnectWebhookURL == "" {
		return nil
	}

	req := &disconnectWebhookRequest{
		RoomID:   c.roomID,
		ClientID: c.ID,
	}

	resp, err := c.postRequest(config.DisconnectWebhookURL, req)
	if err != nil {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Err(err).
			Caller().
			Msg("DiconnectWebhookError")
		return errDisconnectWebhook
	}
	defer resp.Body.Close()

	c.webhookLog("disconnectReq", req)

	// 200 以外で返ってきたときはエラーとする
	if resp.StatusCode != 200 {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Err(err).
			Caller().
			Msg("DiconnectWebhookResponseError")
		return errDisconnectWebhookResponse
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Error().
			Str("roomlId", c.roomID).
			Str("clientId", c.ID).
			Bytes("body", body).
			Err(err).
			Caller().
			Msg("DiconnectWebhookResponseError")
		return err
	}

	// ログ出力用$
	httpResponse := &httpResponse{
		Status: resp.Status,
		Proto:  resp.Proto,
		Header: resp.Header,
		Body:   string(body),
	}

	c.webhookLog("disconnectResp", httpResponse)

	return nil
}
