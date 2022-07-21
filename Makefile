.PHONY: all clean windows mac
### バージョンの定義
VERSION     := "v1.12.1"
COMMIT      := $(shell git rev-parse --short HEAD)

### コマンドの定義
LDFLAGS  = "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

### PHONY ターゲットのビルドルール
all: windows mac

clean:
	rm -rf build/bin/twLauncher*
	rm -rf build/bin/*.exe

windows: build/bin/twsnmpfc-amd64-installer.exe

mac:
	wails build  -platform darwin -ldflags $(LDFLAGS)


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
