package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
)

func TestSignalingHandler(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(signalingHandler))
	defer s.Close()

	// http://127.0.0.1 を ws://127.0.0.1 に変更する
	u := "ws" + strings.TrimPrefix(s.URL, "http")

	ws, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		t.Fatalf("%v", err)
	}
	defer ws.Close()

	// JSON
	if err := ws.WriteMessage(websocket.TextMessage, []byte("abc")); err != nil {
		// assert.Equal(t, http.StatusOK, resp.StatusCode)
		t.Fatalf("%v", err)
	}
	// _, _, err = ws.ReadMessage()
	// if err != nil {
	// 	t.Fatalf("%v", err)
	// }
}
