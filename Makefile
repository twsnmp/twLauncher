.PHONY: all clean windows mac windebug wininstaller dev
### バージョンの定義
VERSION     := "v1.36.3"
COMMIT      := $(shell git rev-parse --short HEAD)

### コマンドの定義
LDFLAGS  = "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

### PHONY ターゲットのビルドルール
all: windows mac

clean:
	rm -rf build/bin/twLauncher*
	rm -rf build/bin/*.exe

windows: build/bin/twLauncher.exe

mac: build/bin/twLauncher.app/Contents/MacOS/twLauncher

build/bin/twLauncher.app/Contents/MacOS/twLauncher:
	wails build  -platform darwin -ldflags $(LDFLAGS)

windebug:
	wails build  -platform windows/amd64 -debug -clean -windowsconsole

dev:
	 wails dev -e go,svelte

build/bin/twLauncher.exe: 
	wails build  -platform windows/amd64  -ldflags $(LDFLAGS)
