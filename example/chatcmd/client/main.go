package main

import (
	"fmt"
	"log"
	"time"

	"github.com/anthony/gopher-pipe/example/chat"
)

func main() {
	cli, err := chat.NewChatClient("127.0.0.1:9200")
	if err != nil {
		log.Fatalln(err)
	}
	defer cli.Close()
	ok, err := cli.Login("alice")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Login ok:", ok)
	// Keep brief wait for demo
	time.Sleep(time.Millisecond * 100)
}
