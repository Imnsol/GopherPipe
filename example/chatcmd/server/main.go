package main

import (
	"fmt"
	"log"

	"github.com/anthony/gopher-pipe/gopherpipe"
)

// impl of chat.ChatService
type chatImpl struct{}

func (c *chatImpl) Login(user string) (bool, error) {
	log.Printf("Login called for: %s", user)
	return true, nil
}

func main() {
	srv := gopherpipe.NewServer(":9200")
	srv.Register("ChatService", &chatImpl{})
	fmt.Println("Chat service listening :9200")
	if err := srv.Serve(); err != nil {
		panic(err)
	}
}
