package main

import (
	"embed"
	"fmt"
	"log"
	"runtime"

	"github.com/wailsapp/wails/v2/pkg/options/mac"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
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
	myLogger := logger.NewDefaultLogger()
	if runtime.GOOS == "windows" {
		myLogger = logger.NewFileLogger("./twl.log")
	}
	// Create application with options
	err := wails.Run(&options.App{
		Title:             "TWSNMP起動/設定ツール",
		Width:             720,
		Height:            600,
		MinWidth:          720,
		MinHeight:         650,
		MaxWidth:          1280,
		MaxHeight:         1000,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		RGBA:              &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		Assets:            assets,
		Logger:            myLogger,
		LogLevel:          logger.DEBUG,
		OnStartup:         app.startup,
		OnDomReady:        app.domReady,
		OnShutdown:        app.shutdown,
		Bind: []interface{}{
			app,
		},
		// Windows platform specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
		},
		Mac: &mac.Options{
			TitleBar: &mac.TitleBar{
				TitlebarAppearsTransparent: true,
				HideTitle:                  true,
				HideTitleBar:               true,
				FullSizeContent:            true,
				UseToolbar:                 false,
				HideToolbarSeparator:       true,
			},
			Appearance:           mac.NSAppearanceNameDarkAqua,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
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
