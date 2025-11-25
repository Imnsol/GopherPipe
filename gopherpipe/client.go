package gopherpipe

import (
	"encoding/gob"
	"fmt"
	"net"
	"sync/atomic"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

type Client struct {
	conn    net.Conn
	counter uint64
}

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

func (c *Client) nextID() uint64 {
	return atomic.AddUint64(&c.counter, 1)
}

// CallUnary sends a unary call and waits for response
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
