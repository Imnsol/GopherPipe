// Package gopherpipe contains a tiny prototype client/server and envelope
// primitives used by the examples in this repository. It's intentionally
// minimal and serves as a PoC for a Go-first RPC approach built on small,
// efficient wire frames.
package gopherpipe

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

// Client is a tiny RPC client used by the example client stubs in this repo.
// It keeps a single TCP connection and a request counter used for Call IDs.
type Client struct {
	conn    net.Conn
	counter uint64
}

// Dial connects to a TCP address and returns a Client ready to send RPCs.
// For the prototype we perform minimal negotiation and register example
// types with gob for encoding/decoding.
func Dial(addr string) (*Client, error) {
	c, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	// Minimal negotiation: skipping for prototype
	// Register gob for Envelope
	gob.Register(Envelope{})
	return &Client{conn: c}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// nextID returns an incremented counter used for unique call identifiers.
func (c *Client) nextID() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

// CallUnary sends a unary call and waits for response
// CallUnary performs a unary RPC: it encodes payload, sends a data frame to
// the server, waits for a response and decodes it into out.
func (c *Client) CallUnary(service, method string, payload interface{}, out interface{}) error {
	b, err := codec.Encode(payload)
	if err != nil {
		return err
	}
	env := Envelope{RPCType: Unary, ServiceName: service, MethodName: method, CallID: c.nextID(), Body: b}
	envb, err := codec.Encode(env)
	if err != nil {
		return err
	}
	if err := tcplite.WriteFrame(c.conn, tcplite.FrameTypeData, envb); err != nil {
		return err
	}
	// wait for response
	ftype, framePayload, err := tcplite.ReadFrame(c.conn)
	if err != nil {
		return err
	}
	if ftype != tcplite.FrameTypeData {
		return fmt.Errorf("unexpected frame: %d", ftype)
	}
	var resp Envelope
	if err := codec.Decode(framePayload, &resp); err != nil {
		return err
	}
	if resp.CallID != env.CallID {
		return fmt.Errorf("mismatched call id")
	}
	// unmarshal response body into out
	if err := codec.Decode(resp.Body, out); err != nil {
		return err
	}
	return nil
}
