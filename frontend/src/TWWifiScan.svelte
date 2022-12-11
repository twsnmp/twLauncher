<script>
  import "../node_modules/98.css/dist/98.css";
  import Alert from "./Alert.svelte";
  export let env = "";
  let info = {
    Params: [],
    Running: false,
    Task: false,
  };
  let conf = {
    Syslog: "",
    Interval: 600,
    Retention: 3600,
  };
  let alert = {
    title: "",
    message: "",
  };
  window.go.main.App.GetProcessInfo("twWifiScan").then((r) => {
    info = r;
    setConf(r.Params);
  });
  const start = () => {
    if (!checkParams()) {
      return;
    }
    let params = [];
    params.push("-syslog");
    params.push(conf.Syslog);
    params.push("-interval");
    params.push(conf.Interval + "");
    setAlert(
      "twWifiScan起動中",
      "twWifiScanを起動しています。お待ちください。",
      true
    );
    window.go.main.App.Start("twWifiScan", params, conf.Task).then((r) => {
      if (r === "") {
        info.Running = true;
        setAlert("", "", false);
      } else {
        setAlert("twWifiScan起動エラー", r, false);
      }
    });
  };
  const stop = () => {
    setAlert(
      "twWifiScan停止中",
      "twWifiScanを停止しています。お待ちください。",
      true
    );
    window.go.main.App.Stop("twWifiScan").then((r) => {
      if (r === "") {
        info.Running = false;
        setAlert("", "", false);
      } else {
        setAlert("twWifiScan停止エラー", r, false);
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
        "twWifiScanパラメータエラー",
        "Syslogの送信先を指定してください。",
        false
      );
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
        case "-interval":
          if (i < params.length - 1) {
            conf.Interval = params[i + 1] * 1;
            i++;
          }
          break;
      }
    }
    if (conf.Interval < 60) {
      conf.Interval = 600;
    }
    conf.Task = info.Task && env == "windows";
  };
</script>

<Alert prop={alert} />
<fieldset>
  <legend>Wifiアクセスポイントセンサー(twWifiScan)</legend>
  <div class="field-row">
    <label for="status">状態</label>
    <input id="status" disabled type="text" value="{info.Running ? '稼働' : '停止'}"/>
  </div>
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
    <label for="interval">送信間隔:</label>
    <input
      id="interval"
      type="number"
      min="60"
      max="3600"
      bind:value={conf.Interval}
    />
    <label for="interval">秒</label>
    {#if env == "windows"}
      <input type="checkbox" id="wifi_task" bind:checked={conf.Task} />
      <label for="wifi_task">スケジューラー</label>
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

