interface go {
  "main": {
    "App": {
		GetDataStore():Promise<string>
		GetInfo():Promise<LanuncherInfo>
		GetProcessInfo(arg1:string):Promise<ProcessInfo>
		Start(arg1:string,arg2:Array<string>,arg3:boolean):Promise<string>
		Stop(arg1:string):Promise<string>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
