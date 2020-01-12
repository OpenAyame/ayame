package main

import "github.com/rs/zerolog"

// room_id と client_id をログに出したい

func (c *client) signalingLog(message message, rawMessage []byte) {
	if message.Type != "pong" {
		signalingLogger.Log().
			Str("roomId", c.roomID).
			Str("clientId", c.ID).
			Str("type", message.Type).
			Msg(string(rawMessage))
	}
}

func (c *client) debugLog() *zerolog.Event {
	return logger.Debug().
		Str("roomId", c.roomID).
		Str("clientId", c.ID)
}

func (c *client) errLog() *zerolog.Event {
	return logger.Error().
		Str("roomId", c.roomID).
		Str("clientId", c.ID)
}

func (c *client) debugSignalingLog(rawMsg string) {
	logger.Debug().
		Str("roomId", c.roomID).
		Str("clientId", c.ID).
		Str("rawMsg", rawMsg).
		Msg("SIGNALING")
}
