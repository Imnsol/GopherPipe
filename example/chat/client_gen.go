package chat

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
	// Create a gob payload of the argument
	// We reuse codec package for simple encode/decode
	var out bool
	if err := cc.c.CallUnary("ChatService", "Login", user, &out); err != nil {
		return false, err
	}
	return out, nil
}
