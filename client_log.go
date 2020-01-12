package main

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

func (c *client) debugLog(msg string) {
	logger.Debug().
		Str("roomId", c.roomID).
		Str("clientId", c.ID).
		Msg(msg)
}

func (c *client) errorLog(msg string) {
	logger.Error().
		Str("roomId", c.roomID).
		Str("clientId", c.ID).
		Msg(msg)
}

func (c *client) debugSignalingLog(rawMsg string) {
	logger.Debug().
		Str("roomId", c.roomID).
		Str("clientId", c.ID).
		Str("rawMsg", rawMsg).
		Msg("SIGNALING")
}
