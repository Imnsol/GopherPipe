package server

import (
	"bufio"
	"net"
	"strings"
	"testing"
	"time"
)

func TestHTTPProbeGets400(t *testing.T) {
	addr := "127.0.0.1:9310"
	// start server
	go func() { _ = Serve(addr) }()
	// allow start
	time.Sleep(100 * time.Millisecond)

	c, err := net.Dial("tcp", addr)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer c.Close()

	req := "GET / HTTP/1.1\r\nHost: localhost\r\n\r\n"
	_, err = c.Write([]byte(req))
	if err != nil {
		t.Fatalf("write: %v", err)
	}

	c.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
	r := bufio.NewReader(c)
	line, err := r.ReadString('\n')
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if !strings.Contains(line, "400") {
		t.Fatalf("expected HTTP 400 status line; got: %q", line)
	}
}
