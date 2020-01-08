package main

// Type を確認する用
type message struct {
	Type string `json:"type"`
}

type registerMessage struct {
	Type          string       `json:"type"`
	RoomID        string       `json:"roomId"`
	ClientID      string       `json:"clientId"`
	AuthnMetadata *interface{} `json:"authnMetadata"`
	SignalingKey  *string      `json:"signalingKey"`
	// 後方互換性対応
	Key *string `json:"key"`
	// Ayame クライアント情報が詰まっている
	AyameClient *string `json:"ayameClient"`
	Libwebrtc   *string `json:"libwebrtc"`
	Environment *string `json:"environment"`
}

type pingMessage struct {
	Type string `json:"type"`
}

type byeMessage struct {
	Type string `json:"type"`
}

// なにか問題があって閉じる時はこれを使う
// type errorMessage struct {
// 	Type   string `json:"type"`
// 	Reason string `json:"reason"`
// }

type acceptMessage struct {
	Type          string       `json:"type"`
	AuthzMetadata *interface{} `json:"authzMetadata,omitempty"`
	IceServers    *[]iceServer `json:"iceServers,omitempty"`
	IsExistClient bool         `json:"isExistClient"`
	// 後方互換性対応
	IsExistUser bool `json:"isUserClient"`
}

type rejectMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}

type iceServer struct {
	Urls       []string `json:"urls"`
	UserName   *string  `json:"username,omitempty"`
	Credential *string  `json:"credential,omitempty"`
}
