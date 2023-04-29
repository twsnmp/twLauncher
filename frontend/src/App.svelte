<script>
  import TWSNMP from "./TWSNMP.svelte";
  import TWPCAP from "./TWPCAP.svelte";
  import TWWifiScan from "./TWWifiScan.svelte";
  import TWWinLog from "./TWWinLog.svelte";
  import {GetInfo} from "../wailsjs/go/main/App";

  let info = {
    Version: "",
    Env: "",
    Ifaces: [],
    PcapVersion: "",
  };
  GetInfo().then((result) => {
    info = result;
  });
</script>

<main>
  <div class="window">
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
