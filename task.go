package main

import (
	"fmt"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type TaskEnt struct {
	Name   string
	Params string
}

// AddTask : タスクスケジューラに登録する
func (b *App) AddTask(name, params string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("AddTask name=%v,params=%s", name, params))
	return ""
}

// DelTask : タスクスケジューラから削除する
func (b *App) DelTask(name string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("DelTask name=%s", name))
	return ""
}
