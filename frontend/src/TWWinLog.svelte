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
    Remote: '',
    User: '',
    Password: '',
  };
  let alert = {
    title: "",
    message: "",
  };
  window.go.main.App.GetProcessInfo("twWinLog").then((r) => {
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
    if( conf.Remote) {
      params.push("-remote");
      params.push(conf.Remote);
      params.push("-user");
      params.push(conf.User);
      params.push("-password");
      params.push(conf.Passowd);
    }
    setAlert("twWinLog起動中", "twWinLogを起動しています。お待ちください。", true);
    window.go.main.App.Start("twWinLog", params,conf.Task).then((r) => {
      if (r === "") {
        info.Running = true;
        setAlert("", "", false);
      } else {
        setAlert("twWinLog起動エラー", r, false);
      }
    });
  };
  const stop = () => {
    setAlert("twWinLog停止中", "twWinLogを停止しています。お待ちください。", true);
    window.go.main.App.Stop("twWinLog").then((r) => {
      if (r === "") {
        info.Running = false;
        setAlert("", "", false);
      } else {
        setAlert("twWinLog停止エラー", r, false);
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
        "twWinLogパラメータエラー",
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
        case "-remote":
          if (i < params.length - 1) {
            conf.Remote = params[i + 1];
            i++;
          }
          break;
        case "-user":
          if (i < params.length - 1) {
            conf.User = params[i + 1];
            i++;
          }
        case "-password":
          if (i < params.length - 1) {
            conf.Password = params[i + 1];
            i++;
          }
          break;
      }
    }
    conf.Task = info.Task && env == "windows";
  };
</script>

<Alert prop={alert} />
<fieldset>
  <legend>Windowsイベントログセンサー(twWinLog)</legend>
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
    <label for="remote">リモート:</label>
    <input id="remote" type="text" bind:value={conf.Remote} />
  </div>
  {#if info.Remote}
    <div class="field-row">
      <label for="user">ユーザー:</label>
      <input id="user" type="text" bind:value={conf.User} />
      <label for="password">パスワード:</label>
      <input id="password" type="password" bind:value={conf.Password} />
    </div>
  {/if}
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
      <input type="checkbox" id="winlog_task" bind:checked={conf.Task} />
      <label for="winlog_task">スケジューラー</label>
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
