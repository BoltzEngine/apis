package gcm

import (
	"testing"
)

func TestProtocolVersion(t *testing.T) {
	tests := []struct {
		s string
		v int
	}{
		{s: "https://fcm.googleapis.com/v1/projects/myproject-b5ae1/messages:send", v: FcmHttpV1Api},
		{s: "https://fcm.googleapis.com/fcm/send", v: FcmLegacyHttpApi},
		{s: "fcm-xmpp.googleapis.com:5235", v: FcmXmppApi},
		{s: "fcm-xmpp.googleapis.com", v: FcmEndpointError},
		{s: "https://gcm-http.googleapis.com/gcm/send", v: GcmEndpoint},
		{s: "gcm-xmpp.googleapis.com:5235", v: GcmEndpoint},
		{s: "gcm-xmpp.googleapis.com", v: FcmEndpointError},
		{s: "http://example.com/v1/messages", v: FcmEndpointError},
		{s: "https://example.com/v1a/messages", v: FcmEndpointError},
		{s: "https://example.com/pv1/messages", v: FcmEndpointError},
		{s: "https://v1/messages", v: FcmEndpointError},
		{s: "example.com:1234", v: FcmEndpointError},
	}
	for _, tt := range tests {
		v := ProtocolVersion(tt.s)
		if v != tt.v {
			t.Errorf("ProtocolVersion(%s) = %d; want %d\n", tt.s, v, tt.v)
		}
	}
}
