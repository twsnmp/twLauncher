package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/shirou/gopsutil/v3/process"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// ProcessInfo : プロセスの情報
type ProcessInfo struct {
	Name    string
	Params  []string
	Running bool
	Task    bool
}

// GetProcessInfoList : プロセス情報のリストを取得する
func (b *App) GetProcessInfoList() []ProcessInfo {
	wails.LogDebug(b.ctx, "GetProcessInfoList")
	ret := []ProcessInfo{}
	for name, params := range b.processMap {
		ret = append(ret, ProcessInfo{
			Name:    name,
			Params:  params,
			Running: b.findProcess(name) != nil,
			Task:    b.findTask(name) != nil,
		})
	}
	return ret
}

func (b *App) DeleteProcess(name string) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("DeleteProcess name=%s", name))
	delete(b.processMap, name)
	return ""
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
	a := strings.SplitN(name, ":", 2)
	if len(a) == 2 {
		ret = a[0]
	}
	if p, err := os.Executable(); err == nil {
		ret = filepath.Join(filepath.Dir(p), ret)
	} else if dir, err := os.Getwd(); err == nil {
		ret = path.Join(dir, ret)
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
	info := b.GetInfo()
	if info.Env == "windows" {
		if b.findTask(name) == nil {
			if strings.HasPrefix(name, "twsnmpfc") {
				for i := 0; i < 15; i++ {
					p := b.findProcess(name)
					if p == nil {
						break
					}
					b.stopByUdp(int(p.Pid))
					time.Sleep(time.Second * 2)
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
	}
	if p := b.findProcess(name); p != nil {
		if err := p.SendSignal(syscall.SIGINT); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
		}
		if err := p.SendSignal(syscall.SIGTERM); err != nil {
			wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
		}
		if runtime.GOOS == "windows" {
			if strings.HasPrefix(name, "twsnmpfc") {
				for i := 0; i < 15; i++ {
					b.stopByUdp(int(p.Pid))
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

// Delete: プロセスを削除する
func (b *App) Delete(name string) string {
	if p := b.findProcess(name); p != nil {
		return "稼働中です。"
	}
	delete(b.processMap, name)
	wails.LogDebugf(b.ctx, "Stop name=%v process not found", name)
	return ""
}

// findProcess : 起動中のプロセスを探す
func (b *App) findProcess(name string) *process.Process {
	params, ok := b.processMap[name]
	if !ok {
		return nil
	}
	a := strings.SplitN(name, ":", 2)
	if len(a) == 2 {
		name = a[0]
	}
	list, err := process.Processes()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("find process list name=%v err=%v", name, err))
		return nil
	}
	for _, p := range list {
		if n, err := p.Name(); err != nil || !strings.HasPrefix(n, name) {
			continue
		} else {
			log.Println(n)
		}
		if cls, err := p.CmdlineSlice(); err == nil && len(cls) > 1 && cmpArgs(params, cls[1:]) {
			log.Println("hit")
			return p
		}
	}
	wails.LogError(b.ctx, fmt.Sprintf("process not found name=%v", name))
	return nil
}

func cmpArgs(s, p []string) bool {
	log.Println(s)
	log.Println(p)
	if len(s) != len(p) {
		return false
	}
	for i := range s {
		if s[i] != p[i] {
			return false
		}
	}
	return true
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
