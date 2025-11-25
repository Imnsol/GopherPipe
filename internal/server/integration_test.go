package server

import (
	"net"
	"testing"
	"time"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/message"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

// startLocalServer starts a server listening on a random local port and returns the address and a stop func.
func startLocalServer(t *testing.T) (addr string, stop func()) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen failed: %v", err)
	}
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		for {
			conn, err := ln.Accept()
			if err != nil {
				return
			}
			go handleConn(conn)
		}
	}()
	return ln.Addr().String(), func() { ln.Close(); <-stopped }
}

// TestEchoServerEndToEnd exercises the server end-to-end by starting a
// temporary server, connecting, sending a message frame and verifying the
// echo response matches the request.
func TestEchoServerEndToEnd(t *testing.T) {
	addr, stop := startLocalServer(t)
	defer stop()

	// wait a little for server to be ready
	time.Sleep(50 * time.Millisecond)

	conn, err := net.DialTimeout("tcp", addr, time.Second)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	msg := message.Message{ID: 999, From: "itest", Body: "hello integration"}
	b, err := codec.Encode(msg)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	if err := tcplite.WriteFrame(conn, tcplite.FrameTypeData, b); err != nil {
		t.Fatalf("write: %v", err)
	}

	// read response
	ftype, payload, err := tcplite.ReadFrame(conn)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if ftype != tcplite.FrameTypeData {
		t.Fatalf("unexpected frame type: %d", ftype)
	}
	var resp message.Message
	if err := codec.Decode(payload, &resp); err != nil {
		t.Fatalf("decode resp: %v", err)
	}
	if resp.ID != msg.ID || resp.From != msg.From || resp.Body != msg.Body {
		t.Fatalf("mismatch resp: %+v vs %+v", resp, msg)
	}
}
