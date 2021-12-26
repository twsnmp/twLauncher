package main

import (
	"context"
	"os/exec"
)

func getCmd(ctx context.Context, path string, params []string) *exec.Cmd {
	return exec.CommandContext(ctx, path, params...)
}
