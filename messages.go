package main

type message struct {
	Type string `json:"type"`
}

type registerMessage struct {
	Type          string       `json:"type"`
	RoomID        string       `json:"roomId"`
	ClientID      string       `json:"clientId"`
	AuthnMetadata *interface{} `json:"authnMetadata,omitempty"`
	SignalingKey  *string      `json:"signalingKey,omitempty"`
	// TODO(nakai): どこかのタイミングで削除する
	Key *string `json:"key,omitempty"`
}

type pingMessage struct {
	Type string `json:"type"`
}

type byeMessage struct {
	Type string `json:"type"`
}

type acceptMessage struct {
	Type       string       `json:"type"`
	Metadata   interface{}  `json:"authzMetadata,omitempty"`
	IceServers *[]iceServer `json:"iceServers,omitempty"`
	// TODO(nakai): IsExsitClient に変更する、ただし下位互換性が壊れるので慎重に
	IsExistUser bool `json:"isExistUser"`
}

type iceServer struct {
	Urls       []string `json:"urls"`
	UserName   *string  `json:"username,omitempty"`
	Credential *string  `json:"credential,omitempty"`
}

type rejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}
