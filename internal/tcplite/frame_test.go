package tcplite

import (
	"bytes"
	"strings"
	"testing"
)

func TestWriteReadFrame(t *testing.T) {
	b := bytes.NewBuffer(nil)
	payload := []byte("hello world")
	if err := WriteFrame(b, FrameTypeData, payload); err != nil {
		t.Fatalf("write: %v", err)
	}
	typeGot, p, err := ReadFrame(b)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if typeGot != FrameTypeData {
		t.Fatalf("expected data frame, got %d", typeGot)
	}
	if string(p) != string(payload) {
		t.Fatalf("payload mismatch: %s", string(p))
	}
}

func TestReadInvalidHeader(t *testing.T) {
	// Simulate an HTTP client or other non-TCP_LITE client that sends ASCII data
	raw := []byte("GET / HTTP/1.1\r\nHost: localhost\r\n\r\n")
	b := bytes.NewBuffer(raw)
	_, _, err := ReadFrame(b)
	if err == nil {
		t.Fatalf("expected error for invalid header, got nil")
	}
	if !strings.Contains(err.Error(), "invalid frame type") {
		t.Fatalf("unexpected error: %v", err)
	}
}
