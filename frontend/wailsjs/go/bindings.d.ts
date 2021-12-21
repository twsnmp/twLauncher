interface go {
  "main": {
    "App": {
		AddTask(arg1:string,arg2:string):Promise<string>
		Close():Promise<string>
		DelTask(arg1:string):Promise<string>
		GetDataStore():Promise<string>
		GetInfo():Promise<TWLanuncherInfo>
		GetTWPCAP():Promise<TWPCAPInfo>
		GetTWSNMP():Promise<TWSNMPInfo>
		GetTWWifiScan():Promise<TWWifiScanInfo>
		GetTWWinLog():Promise<TWWinLogInfo>
		Start(arg1:string,arg2:string):Promise<string>
		Stop(arg1:string):Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
