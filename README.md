# twLauncher
TWSNNMPシリーズの起動、設定GUIツールです。
Wails v2で開発しています。

[![Godoc Reference](https://godoc.org/github.com/twsnmp/twLauncher?status.svg)](http://godoc.org/github.com/twsnmp/twLauncher)
[![Go Report Card](https://goreportcard.com/badge/twsnmp/twLauncher)](https://goreportcard.com/report/twsnmp/twLauncher)

## Overview

TWSNNMPシリーズの起動、設定GUIツールです。
Wails v2で開発しています。

## Status

開発を開始したばかりです。

## Build

開発するためには、wails v2が必要です。

https://wails.io/docs/gettingstarted/installation

に従ってインストールしてください。


ビルドはmakeで行います。
```
$make
```
以下のターゲットが指定できます。
```
  all        全実行ファイルのビルド（省略可能）
  mac        Mac用の実行ファイルのビルド
  windows    Windows用の実行ファイルのビルド
  clean      ビルドした実行ファイルの削除
```

```
$make
```
を実行すれば、Mac OS,Windows用の実行ファイルが、`build/bin`のディレクトリに作成されます。

## Debug

デバックは

```
$wails dev
```

で起動します。

## Run

Mac OS,Windowsの実行ファイルを起動するだけです。


## Copyright

see ./LICENSE

```
Copyright 2021 Masayuki Yamai
```

