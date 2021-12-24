interface go {
  "main": {
    "App": {
		AddTask(arg1:string,arg2:string):Promise<string>
		Close():Promise<string>
		DelTask(arg1:string):Promise<string>
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
