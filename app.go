package main

import (
	"context"
	"fmt"
	"os"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App application struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (b *App) startup(ctx context.Context) {
	// Perform your setup here
	b.ctx = ctx
	wails.LogDebug(b.ctx, "startup")
}

// domReady is called after the front-end dom has been loaded
func (b *App) domReady(ctx context.Context) {
	// Add your action here
	wails.LogDebug(b.ctx, "domReady")
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	// Perform your teardown here
	wails.LogDebug(b.ctx, "shutdown")
}

// Close : プログラムをWindowから終了させる
func (b *App) Close() string {
	wails.LogDebug(b.ctx, "close")
	go func() {
		time.Sleep(time.Millisecond * 10)
		os.Exit(0)
	}()
	return "ok"
}

// GetDataStore : データストアのフォルダを選択する
func (b *App) GetDataStore() string {
	dir, err := wails.OpenDirectoryDialog(b.ctx, wails.OpenDialogOptions{
		Title:                "データストアの選択",
		CanCreateDirectories: true,
	})
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("GetDataStore err=%v", err))
	}
	return dir
}

type TWLanuncherInfo struct {
	Version string
	Env     string
}

// GetInfo : プログラムをWindowから終了させる
func (b *App) GetInfo() TWLanuncherInfo {
	return TWLanuncherInfo{
		Version: fmt.Sprintf("%s(%s)", version, commit),
		// Env:     runtime.GOOS,
		Env: "windows",
	}
}

type TWSNMPInfo struct {
	DataStore string
	Password  string
	Status    string
	Task      string
}

// GetTWSNMP : TWSNMPの情報を取得する
func (b *App) GetTWSNMP() TWSNMPInfo {
	wails.LogDebug(b.ctx, "GetTWSNMP")
	return TWSNMPInfo{
		DataStore: "",
		Password:  "",
		Status:    "稼働中",
		Task:      "No",
	}
}

type TWWinLogInfo struct {
	Syslog   string
	Interval int
	Remote   string
	User     string
	Password string
	Status   string
	Task     string
}

// GetTWWinLog : TWWinLogの情報を取得する
func (b *App) GetTWWinLog() TWWinLogInfo {
	wails.LogDebug(b.ctx, "GetTWSNMP")
	return TWWinLogInfo{
		Password: "",
		Status:   "稼働中",
		Task:     "No",
	}
}
