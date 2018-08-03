package apns

import (
	"net/http"
	"testing"
	"time"
)

func TestInvalidToken(t *testing.T) {
	tab := []struct {
		status Status
		expect bool
	}{
		{status: Success, expect: false},
		{status: ProcessingError, expect: false},
		{status: MissingToken, expect: false},
		{status: MissingTopic, expect: false},
		{status: MissingPayload, expect: false},
		{status: InvalidTokenSize, expect: false},
		{status: InvalidTopicSize, expect: false},
		{status: InvalidPayloadSize, expect: false},
		{status: InvalidToken, expect: true},
		{status: Shutdown, expect: false},
		{status: None, expect: false},
	}
	for _, v := range tab {
		p := &ProtocolError{
			Status: v.status,
		}
		r := p.InvalidToken()
		if r != v.expect {
			t.Errorf("%v: InvalidToken = %v; want %v", v.status, r, v.expect)
		}
	}
}

func TestTimestamp(t *testing.T) {
	var zero time.Time
	now := time.Now()
	tab := []struct {
		Err  ProtocolError
		Time *time.Time // nilはnow以降を表す
	}{
		{
			Err: ProtocolError{Status: InvalidToken},
		},
		{
			Err:  ProtocolError{Status: Shutdown},
			Time: &zero,
		},
		{
			Err:  ProtocolError{StatusCode: http.StatusGone, Time: now},
			Time: &now,
		},
		{
			Err: ProtocolError{StatusCode: http.StatusNotFound},
		},
		{
			Err:  ProtocolError{StatusCode: http.StatusInternalServerError},
			Time: &zero,
		},
	}
	for _, v := range tab {
		tick := v.Err.Timestamp()
		switch {
		case v.Time == nil && !tick.After(now):
			t.Errorf("Timestamp(%+v) = %v; want >%v", v.Err, tick, now)
		case v.Time != nil && !tick.Equal(*v.Time):
			t.Errorf("Timestamp(%+v) = %v; want %v", v.Err, tick, *v.Time)
		}
	}
}

func TestIsLegacyAddr(t *testing.T) {
	tab := []struct {
		Addr   string
		Result bool
	}{
		{Addr: "https://api.push.apple.com/", Result: false},
		{Addr: "https://example.com", Result: false},
		{Addr: "http://example.com:8080/", Result: false},
		{Addr: "gateway.push.apple.com:2195", Result: true},
		{Addr: "example.com:50000", Result: true},
	}
	for _, v := range tab {
		ok := IsLegacyAddr(v.Addr)
		if ok != v.Result {
			t.Errorf("IsLegacyAddr(%q) = %v; want %v", v.Addr, ok, v.Result)
		}
	}
}
