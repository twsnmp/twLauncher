<script>
  import TWSNMP from "./TWSNMP.svelte";
  import TWPCAP from "./TWPCAP.svelte";
  import TWWifiScan from "./TWWifiScan.svelte";
  import TWWinLog from "./TWWinLog.svelte";
  let info = {
    Version: "",
    Env: "",
    Ifaces: [],
    PcapVersion: "",
  };
  window.go.main.App.GetInfo().then((result) => {
    info = result;
  });
</script>

<main>
  <div class="window">
    {#if info.Env == "darwin"}
      <div class="title-bar">
        <div class="title-bar-text">TWSNMP起動/設定ツール</div>
        <div class="title-bar-controls">
          <button aria-label="Minimize" on:click={window.runtime.WindowMinimise} />
          <button aria-label="Close" on:click={window.runtime.Quit} />
        </div>
      </div>
    {/if}
    <div class="window-body" data-wails-no-drag>
      <TWSNMP env={info.Env} />
      <TWPCAP env={info.Env} ifaces={info.Ifaces} pcapVersion={info.PcapVersion} />
      <TWWifiScan env={info.Env} />
      {#if info.Env === "windows" || info.Env === "winStore"}
        <TWWinLog env={info.Env} />
      {/if}
    </div>
    <div class="status-bar">
      <p class="status-bar-field">
        twLauncher {info.Version} Copyright (c) 2021 Masayuki Yamai
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
 :global(fieldset+fieldset) {
    margin-top: 10px;
  }
  :global(label) {
    width: 80px;
    margin-left: 10px;
  }
  :global(input[type="number"]) {
    width: 50px;
  }
  :global(fieldset input[type="text"]) {
    line-height: normal;
  }
  :global(input[type="text"]:disabled) {
    color: black;
    line-height: normal;
  }

</style>
