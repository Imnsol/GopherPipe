package main

import (
	"flag"
	"fmt"
	"os"
)

// This is a toy generator that writes a hard-coded client stub for ChatService
func main() {
	out := flag.String("out", "example/chat/client_gen.go", "path to generated file")
	flag.Parse()
	if err := os.MkdirAll("example/chat", 0755); err != nil {
		fmt.Println("mkdir error:", err)
		os.Exit(1)
	}
	f, err := os.Create(*out)
	if err != nil {
		fmt.Println("create error:", err)
		os.Exit(1)
	}
	defer f.Close()
	f.WriteString(`package chat

import (
	"github.com/anthony/gopher-pipe/gopherpipe"
	"github.com/anthony/gopher-pipe/internal/codec"
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
`)
	fmt.Println("wrote", *out)
}
