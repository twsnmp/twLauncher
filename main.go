package main

import (
	"embed"
	"fmt"
	"log"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var version = "v0.0.0"
var commit = ""

func main() {
	// Create an instance of the app structure
	app := NewApp()
	// Create application with options
	err := wails.Run(&options.App{
		Title:            "TWSNMP起動/設定ツール",
		Width:            1024,
		Height:           768,
		BackgroundColour: &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
		Mac: &mac.Options{
			About: &mac.AboutInfo{
				Title:   "TWSNMP起動/設定ツール",
				Message: fmt.Sprintf("twLauncher %s(%s) © 2021 Masayuki Yamai", version, commit),
				Icon:    icon,
			},
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}
