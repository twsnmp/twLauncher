<script>
	import "../node_modules/98.css/dist/98.css"
	import TWSNMP from "./TWSNMP.svelte"
	import TWPCAP from "./TWPCAP.svelte"
	import TWWifiScan from "./TWWifiScan.svelte"
	import TWWinLog from "./TWWinLog.svelte"
	let info = {
		Version: '',
		Env: '',
	}
	window.go.main.App.GetInfo().then((result) => {
		info = result
	})
	function closeWindow() {
		window.go.main.App.Close().then((result) => {})
	}
	function minWindow() {
		window.runtime.WindowMinimise()
	}
</script>

<main>
	<div class="window">
	{#if info.Env == 'darwin'}
		<div class="title-bar">
			<div class="title-bar-text">TWSNMP起動/設定ツール</div>
			<div class="title-bar-controls">
				<button aria-label="Minimize" on:click={minWindow} />
				<button aria-label="Close" on:click={closeWindow} />
			</div>
		</div>
	{/if}
		<div class="window-body" data-wails-no-drag>
			<TWSNMP/>
			<TWPCAP />
			<TWWifiScan />
			{#if info.Env === 'windows'}
				<TWWinLog />
			{/if}
		</div>
		<div class="status-bar">
			<p class="status-bar-field">
				TWSNMP FC { info.Version } Copyright (c) 2021 Masayuki Yamai
			</p>
		</div>
	</div>
</main>

<style>
	div.window {
		width: 100%;
		height: 100%;
		position: relative;
		padding-bottom: 15px;
		box-sizing: border-box;
		min-height: 100vh;
	}

	div.status-bar {
		width: 99%;
		position: absolute;
		bottom: 0;
	}
</style>
