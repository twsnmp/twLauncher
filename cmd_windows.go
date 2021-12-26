package main

import (
	"context"
	"os/exec"
	"syscall"
)

func getCmd(ctx context.Context, path string, params []string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, path, params...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd
}
