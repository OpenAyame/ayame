package main

type AcceptMessage struct {
	Type       string        `json:"type"`
	IceServers []interface{} `json:"iceServers,omitempty"`
}

type RejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type AcceptMetadataMessage struct {
	Type       string        `json:"type"`
	Metadata   interface{}   `json:"authzMetadata,omitempty"`
	IceServers []interface{} `json:"iceServers,omitempty"`
}
