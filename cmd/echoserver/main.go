package main

// A small echo server binary used for manual testing and examples. It
// delegates to internal/server.Serve to exercise tcplite framing and
// codec handling.

import (
	"fmt"
	"log"
	"os"

	"github.com/anthony/gopher-pipe/internal/server"
)

func main() {
	addr := ":9000"
	if len(os.Args) > 1 {
		addr = os.Args[1]
	}
	fmt.Println("echo server listening on", addr)
	if err := server.Serve(addr); err != nil {
		log.Fatalf("server exited: %v", err)
	}
}
