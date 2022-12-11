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
    DataStore: "",
    Password: "",
    Local: false,
    Task: false,
  };
  let alert = {
    title: "",
    message: "",
  };
  window.go.main.App.GetProcessInfo("twsnmpfc").then((r) => {
    info = r;
    setConf(r.Params);
  });
  const getDataStore = () => {
    window.go.main.App.GetDataStore().then((r) => {
      conf.DataStore = r;
    });
  };
  const start = () => {
    if (!checkParams()) {
      return;
    }
    let params = [];
    params.push("-datastore");
    params.push(conf.DataStore);
    if (conf.Local) {
      params.push("-local");
    }
    if (conf.Password != "") {
      params.push("-password");
      params.push(conf.Password);
    }
    setAlert("TWSNMP起動中", "TWSNMPを起動しています。お待ちください。", true);
    window.go.main.App.Start("twsnmpfc", params,conf.Task).then((r) => {
      if (r === "") {
        info.Running = true;
        setAlert("", "", false);
      } else {
        setAlert("TWSNMP起動エラー", r, false);
      }
    });
  };
  const stop = () => {
    setAlert("TWSNMP停止中", "TWSNMPを停止しています。お待ちください。", true);
    window.go.main.App.Stop("twsnmpfc").then((r) => {
      if (r === "") {
        info.Running = false;
        setAlert("", "", false);
      } else {
        setAlert("TWSNMP停止エラー", r, false);
      }
    });
  };
  const setAlert = (t, m, w) => {
    alert.title = t;
    alert.message = m;
    alert.wait = w;
  };
  const checkParams = () => {
    if (!conf.DataStore) {
      setAlert(
        "TWSNMPパラメータエラー",
        "データストアを選択してください。",
        false
      );
      return false;
    }
    return true;
  };
  const setConf = (params) => {
    conf.Local = false;
    for (let i = 0; i < params.length; i++) {
      switch (params[i]) {
        case "-datastore":
          if (i < params.length - 1) {
            conf.DataStore = params[i + 1];
            i++;
          }
          break;
        case "-password":
          if (i < params.length - 1) {
            conf.Password = params[i + 1];
            i++;
          }
          break;
        case "-local":
          conf.Local = true;
      }
    }
    conf.Task = info.Task && env == "windows";
  };
  const help = () => {
    window.runtime.BrowserOpenURL("https://zenn.dev/twsnmp/books/twsnmpfc-manual");
  }
</script>

<Alert prop={alert} />
<fieldset>
  <legend>TWSNMP FC</legend>
  <div class="field-row">
    <label for="status">状態</label>
    <input id="status" disabled type="text" value="{info.Running ? '稼働' : '停止'}"/>
  </div>
  <div class="field-row">
    <label for="datastore">データストア:</label>
    <input id="datastore" type="text" bind:value={conf.DataStore} />
    <button style="min-width: 20px" on:click={getDataStore}>&lt;</button>
  </div>
  <div class="field-row">
    <label for="twsnmp_password">パスワード :</label>
    <input
      id="twsnmp_password"
      type="password"
      style="margin-right: 10px;"
      bind:value={conf.Password}
    />
    <input type="checkbox" id="local" bind:checked={conf.Local} />
    <label for="local">ローカル</label>
    {#if env == "windows"}
      <input type="checkbox" id="twsnmp_task" bind:checked={conf.Task} />
      <label for="twsnmp_task">スケジューラー</label>
    {/if}
  </div>
  <div class="field-row">
    {#if info.Running}
      <button on:click={stop}>停止</button>
    {:else}
      <button on:click={start}>起動</button>
    {/if}
    <button on:click={help}>ヘルプ</button>
  </div>
</fieldset>

<style>
  #datastore {
    width: 80%;
  }
</style>
