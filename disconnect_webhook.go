package main

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
			Msg("DiconnectWebhookResponseError")
		return errDisconnectWebhookResponse
	}
	defer resp.Body.Close()

	c.webhookLog("disconnectReq", req)
	c.webhookLog("disconnectResp", resp)

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

	return nil
}
