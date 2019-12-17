package main

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
