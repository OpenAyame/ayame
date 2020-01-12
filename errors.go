package main

import (
	"errors"
)

// メッセージはくっつける

var (
	errInvalidMessageType = errors.New("InvalidMessageType")
	errMissingRoomID      = errors.New("MissingRoomID")
	errMissingClientID    = errors.New("MissingClientID")
	errInvalidJSON        = errors.New("InvalidJSON")
	errDuplicateClientID  = errors.New("DuplicateClientID")

	errRegistrationIncomplete = errors.New("RegistrationIncomplete")

	errAuthnWebhook         = errors.New("AuthnWebhookError")
	errAuthnWebhookResponse = errors.New("AuthnWebhookResponseError")
	errAuthnWebhookReject   = errors.New("AuthnWebhookReject")

	errDisconnectWebhook         = errors.New("DisconnectWebhookError")
	errDisconnectWebhookResponse = errors.New("DisconnectWebhookResponseError")

	errRoomFull = errors.New("RoomFull")
	// 想定外のエラー
	// errInternalServer = errors.New("InternalServerError")
)
