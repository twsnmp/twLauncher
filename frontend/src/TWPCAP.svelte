<script>
  import { afterUpdate } from "svelte";
  import "../node_modules/98.css/dist/98.css";
  import Alert from "./Alert.svelte";
  export let env = "";
  export let ifaces = [];
  let info = {
    Params: [],
    Running: false,
    Task: false,
  };
  let conf = {
    Syslog: "",
    Interval: 600,
    Retention: 3600,
    Iface: "",
  };
  let alert = {
    title: "",
    message: "",
  };
  window.go.main.App.GetProcessInfo("twpcap").then((r) => {
    info = r;
    setConf(info.Params);
  });
  afterUpdate(() => {
    if (!conf.Iface) {
      conf.Iface = ifaces && ifaces.length > 0 ? ifaces[0].Value : "";
    }
	});
  const start = () => {
    if (!checkParams()) {
      return;
    }
    let params = [];
    params.push("-syslog");
    params.push(conf.Syslog);
    params.push("-iface");
    params.push(conf.Iface);
    params.push("-interval");
    params.push(conf.Interval + "");
    params.push("-retention");
    params.push(conf.Retention + "");
    setAlert("TWPCAP起動中", "TWPCAPを起動しています。お待ちください。", true);
    window.go.main.App.Start("twpcap", params,conf.Task).then((r) => {
      if (r === "") {
        info.Running = true;
        setAlert("", "", false);
      } else {
        setAlert("TWPCAP起動エラー", r, false);
      }
    });
  };
  const stop = () => {
    setAlert("TWPCAP停止中", "TWPCAPを停止しています。お待ちください。", true);
    window.go.main.App.Stop("twpcap").then((r) => {
      if (r === "") {
        info.Running = false;
        setAlert("", "", false);
      } else {
        setAlert("TWPCAP停止エラー", r, false);
      }
    });
  };
  const setAlert = (t, m, w) => {
    alert.title = t;
    alert.message = m;
    alert.wait = w;
  };
  const checkParams = () => {
    if (!conf.Syslog) {
      setAlert(
        "TWPCAPパラメータエラー",
        "Syslogの送信先を指定してください。",
        false
      );
      return false;
    }
    if (!conf.Iface) {
      setAlert("TWPCAPパラメータエラー", "LAN I/Fを指定してください。", false);
      return false;
    }
    return true;
  };
  const setConf = (params) => {
    for (let i = 0; i < params.length; i++) {
      switch (params[i]) {
        case "-syslog":
          if (i < params.length - 1) {
            conf.Syslog = params[i + 1];
            i++;
          }
          break;
        case "-iface":
          if (i < params.length - 1) {
            conf.Iface = params[i + 1];
            i++;
          }
          break;
        case "-interval":
          if (i < params.length - 1) {
            conf.Interval = params[i + 1] * 1;
            i++;
          }
          break;
        case "-retention":
          if (i < params.length - 1) {
            conf.Retention = params[i + 1] * 1;
            i++;
          }
          break;
      }
    }
    if (conf.Interval < 60) {
      conf.Interval = 600;
    }
    if (conf.Retention < 600) {
      conf.Retention = 3600;
    }
    conf.Task = info.Task;
  };
</script>

<Alert prop={alert} />
<fieldset>
  <div class="field-row">TWPCAP:{info.Running ? "稼働" : "停止"}</div>
  <div class="field-row">
    <label for="syslog">syslog送信先:</label>
    <input
      id="syslog"
      type="text"
      style="width: 80%;"
      bind:value={conf.Syslog}
    />
  </div>
  <div class="field-row">
    <label for="iface">LAN I/F:</label>
    <select id="iface" bind:value={conf.Iface}>
      {#each ifaces as i}
        <option value={i.Value}>{i.Text}</option>
      {/each}
    </select>
  </div>
  <div class="field-row">
    <label for="interval">送信間隔:</label>
    <input
      id="interval"
      type="number"
      min="60"
      max="3600"
      bind:value={conf.Interval}
    />
    <label for="interval">秒</label>
    <label for="retention">保存時間:</label>
    <input id="retention" type="number" min="600" bind:value={conf.Retention} />
    <label for="retention">秒</label>
    {#if env == "windows"}
      <input type="checkbox" id="twpcap_task" bind:checked={conf.Task} />
      <label for="twpcap_task">スケジューラー</label>
    {/if}
  </div>
  <div class="field-row">
    {#if info.Running}
      <button on:click={stop}>停止</button>
    {:else}
      <button on:click={start}>起動</button>
    {/if}
  </div>
</fieldset>

<style>
  label {
    width: 80px;
    margin-left: 10px;
  }
  input[type="number"] {
    width: 50px;
  }
</style>
