package server

import (
	"net"
	"testing"
	"time"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/message"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

func TestServerEcho(t *testing.T) {
	addr := "127.0.0.1:9100"
	// run server in a goroutine
	go func() {
		_ = Serve(addr)
	}()

	// wait for server to start
	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	msg := message.Message{ID: 42, From: "test", Body: "hello"}
	b, err := codec.Encode(msg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := tcplite.WriteFrame(conn, tcplite.FrameTypeData, b); err != nil {
		t.Fatalf("write: %v", err)
	}
	ftype, payload, err := tcplite.ReadFrame(conn)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if ftype != tcplite.FrameTypeData {
		t.Fatalf("unexpected frame: %d", ftype)
	}
	var resp message.Message
	if err := codec.Decode(payload, &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.ID != msg.ID || resp.Body != msg.Body || resp.From != msg.From {
		t.Fatalf("mismatch: got %+v want %+v", resp, msg)
	}
}
