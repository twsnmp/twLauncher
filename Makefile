.PHONY: all clean windows mac
### バージョンの定義
VERSION     := "v1.7.0"
COMMIT      := $(shell git rev-parse --short HEAD)

### コマンドの定義
LDFLAGS  = "-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)"

### PHONY ターゲットのビルドルール
all: windows mac

clean:
	rm -rf build/bin/twLauncher*

windows:
	wails build  -platform windows -ldflags $(LDFLAGS)

mac:
	wails build  -platform darwin -ldflags $(LDFLAGS)

