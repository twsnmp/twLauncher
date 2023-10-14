package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// findTask : タスクを検索する
func (b *App) findTask(name string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not windows")
	}
	tn := strings.ReplaceAll("\\TWSNMP\\"+name, ":", "_")
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/Query", "/TN", tn, "/XML"})
	return cmd.Run()
}

// createTask : タスクを登録する
func (b *App) createTask(name string, params []string) error {
	if runtime.GOOS != "windows" {
		return fmt.Errorf("not windows")
	}
	tn := strings.ReplaceAll("\\TWSNMP\\"+name, ":", "_")
	// 仮のタスクを登録
	cmd := getCmd(b.ctx, "schtasks.exe",
		[]string{"/Create",
			"/TN", tn,
			"/SC", "ONSTART",
			"/NP",
			"/TR", "'" + b.getExec(name) + "'" + makeArg(params)})
	err := cmd.Run()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("create task err=%v", err))
		return err
	}
	// 作成したタスクのXMLで取得する
	cmd = getCmd(b.ctx, "schtasks.exe",
		[]string{"/Query",
			"/TN", tn,
			"/XML",
		})
	o, err := cmd.CombinedOutput()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("query task out=%s err=%v", strings.TrimSpace(string(o)), err))
		return err
	}
	// タスクを削除する
	cmd = getCmd(b.ctx, "schtasks.exe",
		[]string{"/Delete",
			"/TN", tn,
			"/F",
		})
	err = cmd.Run()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("delete task err=%v", err))
		return err
	}
	// 仮タスクのXMLを書き換える
	xml := string(o)
	if strings.Contains(xml, "<Hidden>false</Hidden>") {
		xml = strings.Replace(xml, "<Hidden>false</Hidden>", "<Hidden>true</Hidden>", 1)
	} else {
		xml = strings.Replace(xml, " <Settings>", " <Settings>\n<Hidden>true</Hidden>", 1)
	}
	fp, err := os.CreateTemp("", "twlauncher")
	if err != nil {
		return err
	}
	defer os.Remove(fp.Name())
	tmpName := fp.Name()
	_, err = fp.WriteString(xml)
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("create task write err=%v", err))
		return err
	}
	fp.Close()
	// XMLからタスクを登録する
	cmd = getCmd(b.ctx, "schtasks.exe", []string{"/Create", "/TN", tn, "/XML", tmpName})
	err = cmd.Run()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("create task from xml err=%v", err))
		return err
	}
	// 作成したタスクのXMLを確認したい時は以下をコメントアウト
	// wails.LogDebug(b.ctx, fmt.Sprintf("create task out=%s", string(o)))
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
	tn := strings.ReplaceAll("\\TWSNMP\\"+name, ":", "_")
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/Delete", "/TN", tn, "/F"})
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
	tn := strings.ReplaceAll("\\TWSNMP\\"+name, ":", "_")
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/Run", "/TN", tn, "/I"})
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
	tn := strings.ReplaceAll("\\TWSNMP\\"+name, ":", "_")
	cmd := getCmd(b.ctx, "schtasks.exe", []string{"/End", "/TN", tn})
	o, err := cmd.CombinedOutput()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("end task out=%s err=%v", strings.TrimSpace(string(o)), err))
		return err
	}
	return nil
}
