package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

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
	b.ctx = ctx
	wails.LogDebug(b.ctx, "startup")
	if err := b.loadConfig(); err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("load config err=%v", err))
	}
}

// shutdown is called at application termination
func (b *App) shutdown(ctx context.Context) {
	wails.LogDebug(b.ctx, "shutdown")
	if err := b.saveConfig(); err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("save config err=%v", err))
	}
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
	Version     string
	Env         string
	Ifaces      []selectDataEnt
	PcapVersion string
}

// GetInfo : ランチャー情報の取得
func (b *App) GetInfo() LanuncherInfo {
	env := runtime.GOOS
	if env == "windows" {
		if p, err := os.Executable(); err == nil {
			if strings.Contains(p, "WindowsApps") {
				env = "winStore"
			}
		}
	}
	return LanuncherInfo{
		Version:     fmt.Sprintf("%s(%s)", version, commit),
		Env:         env,
		Ifaces:      b.getIfaces(),
		PcapVersion: b.pcapVersion(),
	}
}

func (b *App) pcapVersion() string {
	_, err := pcap.FindAllDevs()
	if err != nil {
		wails.LogError(b.ctx, fmt.Sprintf("pcapVersion err=%v", err))
		return ""
	}
	wails.LogInfo(b.ctx, pcap.Version())
	return pcap.Version()
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

func (b *App) loadConfig() error {
	conf, err := getConfigName()
	if err != nil {
		return err
	}
	j, err := os.ReadFile(conf)
	if err != nil {
		return err
	}
	wails.LogDebug(b.ctx, string(j))
	b.processMap = make(map[string][]string)
	return json.Unmarshal(j, &b.processMap)
}

func (b *App) saveConfig() error {
	conf, err := getConfigName()
	if err != nil {
		return err
	}
	j, err := json.Marshal(&b.processMap)
	if err != nil {
		return err
	}
	wails.LogDebug(b.ctx, string(j))
	os.WriteFile(conf, j, 0600)
	return nil
}

func getConfigName() (string, error) {
	c, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return path.Join(c, "twlauncher.conf"), nil
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
