package main

import (
	"context"
	"os/exec"
	"testing"
	"time"
)

// TestBuildGenstub verifies the dummy generator under cmd/genstub compiles
// and can be listed by the go tooling. This is a lightweight build check
// to guard example CI.
func TestBuildGenstub(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "list", "github.com/anthony/gopher-pipe/cmd/genstub")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
}
