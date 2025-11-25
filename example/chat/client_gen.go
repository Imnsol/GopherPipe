package chat

// Generated client helper used by the example chat CLI and tests. This file
// is intentionally small and acts as a toy client stub produced by the
// repository's code generator in the examples.

import (
	"github.com/anthony/gopher-pipe/gopherpipe"
)

// ChatClient is a generated client stub (toy generator)
type ChatClient struct {
	c *gopherpipe.Client
}

func NewChatClient(addr string) (*ChatClient, error) {
	c, err := gopherpipe.Dial(addr)
	if err != nil {
		return nil, err
	}
	return &ChatClient{c: c}, nil
}

func (cc *ChatClient) Close() error {
	return cc.c.Close()
}

func (cc *ChatClient) Login(user string) (bool, error) {
	// Login calls the remote ChatService.Login method via the small
	// gopherpipe client. The result is decoded into a bool.
	var out bool
	if err := cc.c.CallUnary("ChatService", "Login", user, &out); err != nil {
		return false, err
	}
	return out, nil
}
