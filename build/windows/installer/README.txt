1.起動ツール(twLauncher)
スタートメニューの「TWSNMP FC」ー「TWSNMP起動ツール」から起動できます。
スケジューラー経由で起動する場合は、「スケジューラー」にチェックしてください。

2.TWSNMP FC
ネットワーク管理ソフトの本体です。
コマンドプロンプトからは、以下の手順で起動きます。

(1)ターミナルを開いてコピーしたディレクトリに移動してください。

>cd twsnmpfc

(2)データストアのディレクトリを作成してください。

>e:
>cd \
>mkdir twsnmpfc

(3)twsnmpfcを起動してください。

>'C:\Program Files\TWSNMP FC\twsnmpfc.exe' -datastore E:\twsnmpfc

(4)ブラウザーからhttp://127.0.0.1:8080にアクセスしてください。
  ユーザーID:パスワードはtwsnmp:twsnmpです。

起動パラメータは以下です。
Usage of C:\Program Files\TWSNMP FC\twsnmpfc.exe:
  -cpuprofile file
    	write cpu profile to file
  -datastore string
    	Path to Data Store directory (default "./datastore")
  -host string
    	Host Name for TLS Cert
  -ip string
    	IP Address for TLS Cert
  -local
    	Local only
  -memprofile file
    	write memory profile to file
  -password string
    	Master Password (default "twsnmpfc!")
  -ping string
    	ping mode icmp or udp
  -port string
    	port (default "8080")
  -restore string
    	Restore DB file name
  -tls
    	Use TLS

3.twWifiScan
無線LANのアクセスポイントの情報をsyslogで送信するセンサープログラムです。
コマンドプロンプトからは、以下の手順で起動きます。

>'C:\Program Files\TWSNMP FC\twWifiScan.exe' -syslog ＜syslogの送信先>

のように起動できます。

起動パラメータは以下です。

>'C:\Program Files\TWSNMP FC\twWifiScan.exe' -h
Usage of C:\Program Files\TWSNMP FC\twWifiScan.exe:
  -cpuprofile file
    	write cpu profile to file
  -iface string
    	monitor interface (default "wlan0")
  -interval int
    	syslog send interval(sec) (default 600)
  -memprofile file
    	write memory profile to file
  -syslog string
    	syslog destnation list

4. twpcap
パケットキャプチャーでネットワーク解析した結果をsyslogで送信するセンサー
プログラムです。
コマンドプロンプトからは、以下の手順で起動きます。

(1)キャプチャできるネットワークインターフェイスを調べる

>'C:\Program Files\TWSNMP FC\twpcap.exe' -list

(2)twpcapを起動する

>'C:\Program Files\TWSNMP FC\twpcap.exe' -iface "<調べたインタ-フェイス名>" -syslog <syslogの送信先>

で起動できます。

コマンドのパラメータは
Usage of C:\Program Files\TWSNMP FC\twpcap.exe:
  -cpuprofile file
    	write cpu profile to file
  -iface string
    	monitor interface
  -interval int
    	syslog send interval(sec) (default 600)
  -list
    	list interface
  -memprofile file
    	write memory profile to file
  -retention int
    	data retention time(sec) (default 3600)
  -syslog string
    	syslog destnation list

5.twWinlog
Windowsのイベントログをsyslogで送信するセンサープログラムです。
コマンドプロンプトからは、以下の手順で起動きます。

>'C:\Program Files\TWSNMP FC\twWinlog.exe' -syslog ＜syslogの送信先>

のように起動できます。

起動パラメータは以下です。

>'C:\Program Files\TWSNMP FC\twWinlog.exe' -h
Usage of C:\Program Files\TWSNMP FC\twWinlog.exe:
  -auth string
        remote authentication:Default|Negotiate|Kerberos|NTLM
  -cpuprofile file
        write cpu profile to file
  -debug
        Debug Mode
  -interval int
        syslog send interval(sec) (default 300)
  -memprofile file
        write memory profile to file
  -password string
        remote user's password
  -remote string
        remote windows pc
  -syslog string
        syslog destnation list
  -user string
        remote user name

6.詳しい説明
noteの以下のマガジンを参照してください。

https://note.com/twsnmp/m/meed0d0ddab5e

タスクスケジューラーへの登録については、以下の記事も参照してください。

https://note.com/twsnmp/n/nf48d231322a0?magazine_key=meed0d0ddab5e
