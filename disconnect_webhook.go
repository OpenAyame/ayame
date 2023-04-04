package ayame

import (
	"fmt"
	"io"
	"net/url"
	"time"
)

type disconnectWebhookRequest struct {
	RoomID       string `json:"roomId"`
	ClientID     string `json:"clientId"`
	ConnectionID string `json:"connectionId"`
}

func (c *connection) disconnectWebhook() error {
	if c.config.DisconnectWebhookURL == "" {
		return nil
	}

	req := &disconnectWebhookRequest{
		RoomID:       c.roomID,
		ClientID:     c.clientID,
		ConnectionID: c.ID,
	}

	start := time.Now()

	resp, err := c.postRequest(c.config.DisconnectWebhookURL, req)
	if err != nil {
		c.errLog().Err(err).Caller().Msg("DiconnectWebhookError")
		return errDisconnectWebhook
	}
	defer resp.Body.Close()

	c.webhookLog("disconnectReq", req)

	u, err := url.Parse(c.config.DisconnectWebhookURL)
	if err != nil {
		c.errLog().Err(err).Caller().Msg("DisconnectWebhookError")
		return errDisconnectWebhook
	}
	statusCode := fmt.Sprintf("%d", resp.StatusCode)
	m := c.metrics
	m.IncWebhookReqCnt(statusCode, "POST", u.Host, u.Path)
	m.ObserveWebhookReqDur(statusCode, "POST", u.Host, u.Path, time.Since(start).Seconds())
	// TODO: ヘッダーのサイズも計測する
	m.ObserveWebhookResSz(statusCode, "POST", u.Host, u.Path, resp.ContentLength)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.errLog().Bytes("body", body).Err(err).Caller().Msg("DiconnectWebhookResponseError")
		return errDisconnectWebhookResponse
	}

	// ログ出力用
	httpResponse := &httpResponse{
		Status: resp.Status,
		Proto:  resp.Proto,
		Header: resp.Header,
		Body:   string(body),
	}

	// 200 以外で返ってきたときはエラーとする
	if resp.StatusCode != 200 {
		c.errLog().Interface("resp", httpResponse).Caller().Msg("DisconnectWebhookUnexpectedStatusCode")
		return errDisconnectWebhookUnexpectedStatusCode
	}

	c.webhookLog("disconnectResp", httpResponse)

	return nil
}
