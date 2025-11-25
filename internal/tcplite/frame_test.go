package tcplite

import (
	"bytes"
	"strings"
	"testing"
)

// TestWriteReadFrame verifies basic round-trip of WriteFrame and ReadFrame
// for a simple payload and ensures frame type and payload are preserved.
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

// TestReadInvalidHeader ensures a non-TCP_LITE header (e.g. an HTTP
// request) results in an InvalidFrameHeaderError so callers can detect and
// handle the situation.
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
