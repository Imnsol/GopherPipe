package main

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

// TestBuildEchoClient is a simple build-time check for the echo client
// command; it uses `go list` to confirm the package is buildable in CI.
func TestBuildEchoClient(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "list", "github.com/anthony/gopher-pipe/cmd/echoclient")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
}
