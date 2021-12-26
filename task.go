package main

import (
	"fmt"
	"runtime"
	"strings"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// findTask : タスクを検索する
func (b *App) findTask(name string) error {
	if runtime.GOOS != "windows" {
		return nil
	}
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/Query", "/TN", "\\TWSNMP\\" + name, "/XML"})
	return cmd.Run()
}

// createTask : タスクを登録する
func (b *App) createTask(name string, params []string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not windows")
	}
	cmd := getCmd(b.ctx, "schtasks.exe",
		[]string{"/Create",
			"/TN", "\\TWSNMP\\" + name,
			"/SC", "ONSTART",
			"/TR", "'" + b.getExec(name) + "'" + makeArg(params)})
	o, err := cmd.CombinedOutput()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("create task out=%s err=%v", strings.TrimSpace(string(o)), err))
		return err
	}
	return nil
}
func makeArg(params []string) string {
	ret := ""
	for _, p := range params {
		ret += " "
		if strings.HasPrefix(p, "-") {
			ret += p
		} else {
			ret += "\"" + p + "\""
		}
	}
	return ret
}

// deleteTask : タスクを削除する
func (b *App) deleteTask(name string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not windows")
	}
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/Delete", "/TN", "\\TWSNMP\\" + name, "/F"})
	o, err := cmd.CombinedOutput()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("delete task out=%s err=%v", strings.TrimSpace(string(o)), err))
		return err
	}

	return nil
}

// runTask : タスクを起動する
func (b *App) runTask(name string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not windows")
	}
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/Run", "/TN", "\\TWSNMP\\" + name, "/I"})
	o, err := cmd.CombinedOutput()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("run task out=%s err=%v", strings.TrimSpace(string(o)), err))
		return err
	}
	return nil
}

// endTask : タスクを停止する
func (b *App) endTask(name string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not windows")
	}
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/End", "/TN", "\\TWSNMP\\" + name})
	o, err := cmd.CombinedOutput()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("end task out=%s err=%v", strings.TrimSpace(string(o)), err))
		return err
	}
	return nil
}
