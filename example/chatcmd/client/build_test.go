package main

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

func TestBuildExampleChatClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "list", "github.com/anthony/gopher-pipe/example/chatcmd/client")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
}
