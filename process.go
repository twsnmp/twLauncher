package main

import (
	"fmt"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// Start : プロセスを起動する
func (b *App) Start(name, params string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Start name=%s,params=%s", name, params))
	return ""
}

// Stop: プロセスを停止する
func (b *App) Stop(name string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Stop name=%v", name))
	return ""
}
