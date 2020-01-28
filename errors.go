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

	errRegistrationIncomplete = errors.New("RegistrationIncomplete")

	errAuthnWebhook                     = errors.New("AuthnWebhookError")
	errAuthnWebhookResponse             = errors.New("AuthnWebhookResponseError")
	errAuthnWebhookUnexpectedStatusCode = errors.New("AuthnWebhookUnexpectedStatusCode")
	errAuthnWebhookReject               = errors.New("AuthnWebhookReject")

	errDisconnectWebhook                     = errors.New("DisconnectWebhookError")
	errDisconnectWebhookResponse             = errors.New("DisconnectWebhookResponseError")
	errDisconnectWebhookUnexpectedStatusCode = errors.New("DisconnectWebhookUnexpectedStatusCode")

	errConfigInvalidLogLevel = errors.New("ConfigInvalidLogLevel")

	errRoomFull = errors.New("RoomFull")
	// 想定外のエラー
	errInternalServer = errors.New("InternalServerError")
)
