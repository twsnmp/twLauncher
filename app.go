package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/gopacket/pcap"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

// App application struct
type App struct {
	ctx        context.Context
	processMap map[string][]string
}

// NewApp creates a new App application struct
func NewApp() *App {
	pmap := make(map[string][]string)
	return &App{
		processMap: pmap,
	}
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

// selectDataEnt : 選択データー
type selectDataEnt struct {
	Text  string
	Value string
}

// LanuncherInfo : ランチャーの情報
type LanuncherInfo struct {
	Version string
	Env     string
	Ifaces  []selectDataEnt
}

// GetInfo : ランチャー情報の取得
func (b *App) GetInfo() LanuncherInfo {
	return LanuncherInfo{
		Version: fmt.Sprintf("%s(%s)", version, commit),
		// Env:     runtime.GOOS,
		Env:    "windows",
		Ifaces: b.getIfaces(),
	}
}

func (b *App) getIfaces() []selectDataEnt {
	ret := []selectDataEnt{}
	// Find all ifaces
	ifaces, err := pcap.FindAllDevs()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("getIfaces err=%v", err))
		return ret
	}
	for _, i := range ifaces {
		if len(i.Addresses) < 1 {
			continue
		}

		ret = append(ret, selectDataEnt{
			Text:  fmt.Sprintf("%s(%s)", i.Name, getIPv4(i.Addresses)),
			Value: i.Name,
		})
	}
	return ret
}

// getIPv4 : LAN I/FのリストからIPv4アドレスを取得する
func getIPv4(addrs []pcap.InterfaceAddress) string {
	r := ""
	for _, ip := range addrs {
		if ip.IP.To4() != nil {
			return ip.IP.String()
		}
		// IPv4のアドレスがなければ最初のアドレスにする
		if r == "" {
			r = ip.IP.String()
		}
	}
	return r
}
