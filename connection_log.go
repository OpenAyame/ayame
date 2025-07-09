package ayame

import (
	"log/slog"
	"slices"
)

func (c *connection) signalingLog(message message, rawMessage []byte) {
	// signaling の type を指定してフィルターする
	if !slices.Contains(c.config.SignalingLogFilters, message.Type) {
		return
	}

	if c.config.Debug {
		c.signalingLogger.Debug("signaling",
			"roomId", c.roomID,
			"clientId", c.clientID,
			"connectionId", c.ID,
			"type", message.Type,
			"rawMessage", string(rawMessage),
		)
		return
	}

	c.signalingLogger.Info("signaling",
		"roomId", c.roomID,
		"clientId", c.clientID,
		"connectionId", c.ID,
		"type", message.Type,
	)
}

func (c *connection) errLog() *slog.Logger {
	return slog.Default().With(
		"roomId", c.roomID,
		"clientID", c.clientID,
		"connectionId", c.ID,
	)
}

func (c *connection) debugLog() *slog.Logger {
	return slog.Default().With(
		"roomId", c.roomID,
		"clientID", c.clientID,
		"connectionId", c.ID,
	)
}