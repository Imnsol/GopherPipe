package main

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestBuildExampleChatServer(t *testing.T) {
	// TestBuildExampleChatServer is a small build-time check that ensures the example
	// chat server compiles and is discoverable to the go tooling.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "list", "github.com/anthony/gopher-pipe/example/chatcmd/server")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
}
