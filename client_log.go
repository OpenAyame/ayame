package main

import "github.com/rs/zerolog"

func (c *connection) signalingLog(message message, rawMessage []byte) {
	if message.Type != "pong" {
		signalingLogger.Log().
			Str("roomId", c.roomID).
			Str("clientID", c.clientID).
			Str("connectionId", c.ID).
			Str("type", message.Type).
			Msg(string(rawMessage))
	}
}

func (c *connection) errLog() *zerolog.Event {
	return logger.Error().
		Str("roomId", c.roomID).
		Str("clientID", c.clientID).
		Str("connectionId", c.ID)
}

func (c *connection) debugLog() *zerolog.Event {
	return logger.Debug().
		Str("roomId", c.roomID).
		Str("clientID", c.clientID).
		Str("connectionId", c.ID)
}
