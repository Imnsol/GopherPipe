package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/anthony/gopher-pipe/internal/codec"
	"github.com/anthony/gopher-pipe/internal/message"
	"github.com/anthony/gopher-pipe/internal/tcplite"
)

func main() {
	addr := "127.0.0.1:9000"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	msg := message.Message{ID: time.Now().UnixNano(), From: "client", Body: "hello from client"}
	b, err := codec.Encode(msg)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}
	if err := tcplite.WriteFrame(conn, tcplite.FrameTypeData, b); err != nil {
		log.Fatalf("write: %v", err)
	}
	ftype, payload, err := tcplite.ReadFrame(conn)
	if err != nil {
		log.Fatalf("read: %v", err)
	}
	if ftype != tcplite.FrameTypeData {
		log.Fatalf("unexpected frame type: %d", ftype)
	}
	var resp message.Message
	if err := codec.Decode(payload, &resp); err != nil {
		log.Fatalf("decode resp: %v", err)
	}
	fmt.Printf("got echo: %+v\n", resp)
}
