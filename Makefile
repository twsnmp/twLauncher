.PHONY: all clean windows mac windebug wininstaller dev
### バージョンの定義
VERSION     := "v1.17.0"
COMMIT      := $(shell git rev-parse --short HEAD)

### コマンドの定義
LDFLAGS  = "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

### PHONY ターゲットのビルドルール
all: windows mac

clean:
	rm -rf build/bin/twLauncher*
	rm -rf build/bin/*.exe

windows: build/bin/twLauncher.exe

mac:
	wails build  -platform darwin -ldflags $(LDFLAGS)

wininstaller: build/bin/twsnmpfc-amd64-installer.exe

windebug:
	wails build  -platform windows -debug

dev:
	 wails dev -e go,svelte

build/bin/twsnmpfc-amd64-installer.exe: build/bin/twsnmpfc.exe build/bin/twpcap.exe build/bin/twWifiScan.exe build/bin/twWinlog.exe build/bin/twLauncher.exe
	wails build  -platform windows/amd64 -ldflags $(LDFLAGS) -nsis

build/bin/twLauncher.exe: 
	wails build  -platform windows/amd64 -ldflags $(LDFLAGS)

build/bin/twsnmpfc.exe: ../twsnmpfc/dist/twsnmpfc.exe
	cp ../twsnmpfc/dist/twsnmpfc.exe build/bin/

build/bin/twpcap.exe: ../twpcap/dist/twpcap.exe
	cp ../twpcap/dist/twpcap.exe build/bin/

build/bin/twWifiScan.exe: ../twWifiScan/dist/twWifiScan.exe
	cp ../twWifiScan/dist/twWifiScan.exe build/bin/

build/bin/twWinlog.exe: ../twWinlog/dist/twWinlog.exe
	cp ../twWinlog/dist/twWinlog.exe build/bin/
