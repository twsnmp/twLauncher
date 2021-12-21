package main

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/google/gopacket/pcap"
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

type selectDataEnt struct {
	Text  string
	Value string
}

type TWLanuncherInfo struct {
	Version string
	Env     string
}

// GetInfo : プログラムをWindowから終了させる
func (b *App) GetInfo() TWLanuncherInfo {
	return TWLanuncherInfo{
		Version: fmt.Sprintf("%s(%s)", version, commit),
		Env:     runtime.GOOS,
	}
}

// TWSNMPInfo : TWSNMP FCの情報
type TWSNMPInfo struct {
	DataStore string
	Password  string
	Local     bool
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

// TWWinLogInfo : twWinLogセンサーの情報
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
		Interval: 300,
		Password: "",
		Status:   "稼働中",
		Task:     "No",
	}
}

// TWPCAPInfo : TWPCAPセンサーの情報
type TWPCAPInfo struct {
	Syslog    string
	Interval  int
	Retention int
	Iface     string
	Ifaces    []selectDataEnt
	Status    string
	Task      string
}

// GetTWPCAP : TWPCAPの情報を取得する
func (b *App) GetTWPCAP() TWPCAPInfo {
	wails.LogDebug(b.ctx, "GetTWPCAP")
	return TWPCAPInfo{
		Interval: 300,
		Ifaces:   b.getIfaces(),
		Status:   "稼働中",
		Task:     "No",
	}
}

// TWWifiScanInfo : TWWifiScanセンサーの情報
type TWWifiScanInfo struct {
	Syslog   string
	Interval int
	Iface    string
	Ifaces   []selectDataEnt
	Status   string
	Task     string
}

// GetTWWifiScan : TWWifiScanの情報を取得する
func (b *App) GetTWWifiScan() TWWifiScanInfo {
	wails.LogDebug(b.ctx, "GetTWWifiScan")
	return TWWifiScanInfo{
		Interval: 300,
		Ifaces:   b.getIfaces(),
		Status:   "稼働中",
		Task:     "No",
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

// getIPv4 : LAN I/FのリストからIPｖ４アドレスを取得する
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
