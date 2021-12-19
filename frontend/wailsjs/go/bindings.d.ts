interface go {
  "main": {
    "App": {
		Close():Promise<string>
		GetDataStore():Promise<string>
		GetInfo():Promise<TWLanuncherInfo>
    },
  }

}

declare global {
	interface Window {
		go: go;
	}
}
