package main

type AcceptMessage struct {
	Type        string        `json:"type"`
	Metadata    interface{}   `json:"authzMetadata,omitempty"`
	IceServers  []interface{} `json:"iceServers,omitempty"`
	IsExistUser bool          `json:"isExistUser"`
}

type RejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}
