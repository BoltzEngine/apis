package gcm

import (
	"testing"
)

func TestProtocolVersion(t *testing.T) {
	tests := []struct {
		s string
		v int
	}{
		{s: "https://fcm.googleapis.com/v1/projects/myproject-b5ae1/messages:send", v: 1},
		{s: "http://example.com/v1/messages", v: 1},
		{s: "https://example.com/v1a/messages", v: 0},
		{s: "https://example.com/pv1/messages", v: 0},
		{s: "https://v1/messages", v: 0},
		{s: "example.com:1234", v: -1},
	}
	for _, tt := range tests {
		v := ProtocolVersion(tt.s)
		if v != tt.v {
			t.Errorf("ProtocolVersion(%s) = %d; want %d\n", tt.s, v, tt.v)
		}
	}
}
