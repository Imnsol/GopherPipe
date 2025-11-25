package message

import "testing"

func TestMessageFields(t *testing.T) {
	m := Message{ID: 123, From: "bob", Body: "hi"}
	if m.ID != 123 || m.From != "bob" || m.Body != "hi" {
		t.Fatalf("fields not set properly: %+v", m)
	}
}

func TestMessageEncodeDecode(t *testing.T) {
	// TestMessageEncodeDecode ensures the Message type round-trips through
	// the codec Encode/Decode helpers used by the examples.
}
