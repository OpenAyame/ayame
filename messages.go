package main

type acceptMessage struct {
	Type        string        `json:"type"`
	Metadata    interface{}   `json:"authzMetadata,omitempty"`
	IceServers  []interface{} `json:"iceServers,omitempty"`
	IsExistUser bool          `json:"isExistUser"`
}

type rejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}
