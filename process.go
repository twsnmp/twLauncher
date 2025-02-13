package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
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
		task := false
		running := false
		if strings.HasPrefix(name, "http") {
			running = b.checkUrl(name)
		} else {
			running = b.findProcess(name) != nil
			task = b.findTask(name) == nil
		}
		ret = append(ret, ProcessInfo{
			Name:    name,
			Params:  params,
			Running: running,
			Task:    task,
		})
	}
	return ret
}

func (b *App) checkUrl(url string) bool {
	ctx, cancel := context.WithTimeout(b.ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false
	}
	req = req.WithContext(ctx)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	defer res.Body.Close()
	return res.StatusCode < 400
}

// Start : プロセスを起動する
func (b *App) Start(name string, params []string, task bool) string {
	wails.LogDebug(b.ctx, fmt.Sprintf("Start name=%s,params=%v", name, params))
	b.processMap[name] = params
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
	return ""
}

func (b *App) Save(name string) {
	b.processMap[name] = []string{}
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
	if runtime.GOOS == "windows" {
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
				return fmt.Sprintf("タスクを停止できません:%v", err)
			}
			if err := b.deleteTask(name); err != nil {
				wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
				return fmt.Sprintf("タスクを削除できません:%v", err)
			}
			if p := b.findProcess(name); p != nil {
				wails.LogError(b.ctx, fmt.Sprintf("kill name=%v", name))
				p.Kill()
			}
			return ""
		}
	}
	if p := b.findProcess(name); p != nil {
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
		} else {
			if err := p.SendSignal(syscall.SIGINT); err != nil {
				wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
			}
			if err := p.SendSignal(syscall.SIGTERM); err != nil {
				wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v err=%v", name, err))
			}
			for i := 0; i < 5; i++ {
				time.Sleep(time.Second * 2)
				p := b.findProcess(name)
				if p == nil {
					return ""
				}
			}
			wails.LogError(b.ctx, fmt.Sprintf("kill name=%v", name))
			p.Kill()
			return "強制終了しました"
		}
	}
	wails.LogError(b.ctx, fmt.Sprintf("Stop name=%v process not found", name))
	return ""
}

// Delete: プロセスを削除する
func (b *App) Delete(name string) string {
	if !strings.HasPrefix(name, "http") {
		if p := b.findProcess(name); p != nil {
			return "稼働中です。"
		}
	}
	delete(b.processMap, name)
	wails.LogDebugf(b.ctx, "Delete name=%v", name)
	return ""
}

// findProcess : 起動中のプロセスを探す
func (b *App) findProcess(name string) *process.Process {
	_, ok := b.processMap[name]
	if !ok {
		wails.LogError(b.ctx, fmt.Sprintf("processMap not found name=%v", name))
		return nil
	}
	oname := name
	a := strings.SplitN(name, ":", 2)
	cp := ""
	if len(a) == 2 {
		name = a[0]
		cp = strings.TrimSpace(a[1])
	}
	name = strings.ToLower(name)
	list, err := process.Processes()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("find process list name=%v err=%v", name, err))
		return nil
	}
	for _, p := range list {
		if n, err := p.Name(); err != nil || !strings.HasPrefix(strings.ToLower(n), name) {
			continue
		}
		if cls, err := p.CmdlineSlice(); err == nil && checkParam(name, cp, cls) {
			return p
		} else {
			wails.LogErrorf(b.ctx, "checkParam name=%s err=%v cp='%s' cls=%v", name, err, cp, cls)
		}
	}
	wails.LogError(b.ctx, fmt.Sprintf("process not found name=%v", oname))
	return nil
}

func checkParam(name, cp string, p []string) bool {
	if name == "twwinlog" && cp == "" {
		for _, v := range p {
			if v == "-remote" {
				return false
			}
		}
		return true
	}
	for i, v := range p {
		if i+1 < len(p) {
			switch v {
			case "-port", "-remote", "-iface", "--apiPort":
				s := strings.ReplaceAll(p[i+1], "\"", "")
				if s == cp {
					return true
				}
			}
		}
	}
	return false
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

func (b *App) needWindowsPrivilege() bool {
	if runtime.GOOS != "windows" {
		return false
	}
	if list, err := process.Processes(); err == nil {
		for _, p := range list {
			if _, err := p.Name(); err != nil && strings.Contains(err.Error(), "Access is denied.") {
				//  Access is denied.で名前の取得できないプロセスがあるのは、権限がないため
				wails.LogError(b.ctx, fmt.Sprintf("NeedWinPrivilege err=%v", err))
				return true
			}
		}
	} else {
		if strings.Contains(err.Error(), "Access is denied.") {
			wails.LogError(b.ctx, fmt.Sprintf("NeedWinPrivilege err=%v", err))
			return true
		}
	}
	return false
}
