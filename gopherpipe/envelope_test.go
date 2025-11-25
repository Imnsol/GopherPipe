package gopherpipe

import (
	"testing"

	"github.com/anthony/gopher-pipe/internal/codec"
)

// TestEnvelopeEncodeDecode verifies that an Envelope can be encoded and
// decoded using the repo's codec implementation without loss of metadata.
func TestEnvelopeEncodeDecode(t *testing.T) {
	env := Envelope{RPCType: Unary, ServiceName: "S", MethodName: "M", CallID: 1, Body: []byte("payload")}
	b, err := codec.Encode(env)
	if err != nil {
		t.Fatalf("encode env: %v", err)
	}
	var got Envelope
	if err := codec.Decode(b, &got); err != nil {
		t.Fatalf("decode env: %v", err)
	}
	if got.CallID != env.CallID || got.ServiceName != env.ServiceName || got.MethodName != env.MethodName {
		t.Fatalf("mismatch envelope: %+v vs %+v", got, env)
	}
}
