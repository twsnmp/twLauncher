package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/mitchellh/go-ps"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ProcessInfo : プロセスの情報
type ProcessInfo struct {
	Params  []string
	Running bool
	Task    bool
}

// GetProcessInfo : プロセスの情報を取得する
func (b *App) GetProcessInfo(name string) ProcessInfo {
	wails.LogDebug(b.ctx, fmt.Sprintf("GetProcessInfo name=%s", name))
	params, ok := b.processMap[name]
	if !ok {
		params = []string{}
	}
	return ProcessInfo{
		Running: b.findProcess(name) != nil,
		Params:  params,
		Task:    false,
	}
}

// Start : プロセスを起動する
func (b *App) Start(name string, params []string, task bool) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Start name=%s,params=%v", name, params))
	if p := b.findProcess(name); p != nil {
		wails.LogDebug(b.ctx, name+" is running")
		return ""
	}
	cmd := exec.Command(b.getExec(name), params...)
	if err := cmd.Start(); err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("Start name=%v err=%v", name, err))
		return fmt.Sprintf("起動できません err=%v", err)
	}
	if cmd.Process != nil {
		wails.LogDebug(b.ctx, "cmd.Process.Release")
		if err := cmd.Process.Release(); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Start name=%v err=%v", name, err))
		}
	}
	b.processMap[name] = params
	return ""
}

func (b *App) getExec(name string) string {
	ret := name
	if p, err := os.Executable(); err == nil {
		ret = path.Join(path.Dir(p), name)
	} else if dir, err := os.Getwd(); err == nil {
		ret = path.Join(dir, name)
	} else {
		wails.LogError(b.ctx, fmt.Sprintf("getExec name=%v err=%v", name, err))
	}
	if runtime.GOOS == "darwin" {
		ret += ".app"
	} else {
		ret += ".exe"
	}
	return ret
}

// Stop: プロセスを停止する
func (b *App) Stop(name string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Stop name=%v", name))
	if p := b.findProcess(name); p != nil {
		p.Signal(os.Interrupt)
		for i := 0; i < 60; i++ {
			time.Sleep(time.Second * 2)
			p := b.findProcess(name)
			if p == nil {
				return ""
			}
		}
		p.Kill()
		return "強制終了しました"
	}
	wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v process not found", name))
	return ""
}

// findProcess : 起動中のプロセスを探す
func (b *App) findProcess(name string) *os.Process {
	list, err := ps.Processes()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("find process list name=%v err=%v", name, err))
		return nil
	}
	for _, p := range list {
		if strings.Contains(p.Executable(), name) {
			process, err := os.FindProcess(p.Pid())
			if err == nil {
				return process
			}
			wails.LogError(b.ctx, fmt.Sprintf("find process os name=%v err=%v", name, err))
		}
	}
	return nil
}

// findTask : タスクスケジューラのタスクを検索する
