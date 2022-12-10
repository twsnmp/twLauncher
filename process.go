package main

import (
	"fmt"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
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
		Task:    b.findTask(name) == nil,
	}
}

// Start : プロセスを起動する
func (b *App) Start(name string, params []string, task bool) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Start name=%s,params=%v", name, params))
	if task {
		if err := b.endTask(name); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("end task name=%v err=%v", name, err))
		}
		if err := b.deleteTask(name); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("delete task name=%v err=%v", name, err))
		}
		if err := b.createTask(name, params); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("create task name=%v err=%v", name, err))
			return fmt.Sprintf("タスクを登録できません:%v", err)
		}
		if err := b.runTask(name); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("run task name=%v err=%v", name, err))
			return fmt.Sprintf("タスクを起動できません:%v", err)
		}
		b.processMap[name] = params
		return ""
	}
	if p := b.findProcess(name); p != nil {
		wails.LogDebug(b.ctx, name+" is running")
		return ""
	}
	cmd := getCmd(b.ctx, b.getExec(name), params)
	if err := cmd.Start(); err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("Start name=%v err=%v", name, err))
		return fmt.Sprintf("起動できません err=%v", err)
	}
	if cmd.Process != nil {
		go cmd.Process.Wait()
	}
	b.processMap[name] = params
	return ""
}

func (b *App) getExec(name string) string {
	ret := name
	if p, err := os.Executable(); err == nil {
		ret = filepath.Join(filepath.Dir(p), name)
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
	wails.LogDebug(b.ctx, fmt.Sprintf("get exec ret=%s", ret))
	return ret
}

// Stop: プロセスを停止する
func (b *App) Stop(name string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Stop name=%v", name))
	if b.findTask(name) == nil {
		if runtime.GOOS == "windows" {
			if name == "twsnmpfc" {
				for i := 0; i < 15; i++ {
					p := b.findProcess(name)
					if p == nil {
						break
					}
					b.stopByUdp(p.Pid)
					time.Sleep(time.Second * 2)
				}
			}
		}
		if err := b.endTask(name); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
		}
		if err := b.deleteTask(name); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
			return fmt.Sprintf("タスク削除エラー:%v", err)
		}
		return ""
	}
	if p := b.findProcess(name); p != nil {
		if err := p.Signal(syscall.SIGINT); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
		}
		if err := p.Signal(syscall.SIGTERM); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
		}
		if runtime.GOOS == "windows" {
			if name == "twsnmpfc" {
				for i := 0; i < 15; i++ {
					b.stopByUdp(p.Pid)
					time.Sleep(time.Second * 2)
					p := b.findProcess(name)
					if p == nil {
						return ""
					}
				}
			} else {
				p.Kill()
				return ""
			}
		}
		for i := 0; i < 5; i++ {
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

func (b *App) stopByUdp(pid int) {
	for i := 8080; i < 8180; i++ {
		func() {
			conn, err := net.Dial("udp4", fmt.Sprintf("127.0.0.1:%d", i))
			if err == nil {
				defer conn.Close()
				_, err = conn.Write([]byte(fmt.Sprintf("%d", pid)))
				if err != nil {
					wails.LogError(b.ctx, fmt.Sprintf("stopByUdp pid=%v err=%v", pid, err))
				}
			}
		}()
	}
}
