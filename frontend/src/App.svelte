<script>
  import { onMount } from "svelte";
  import logo from "./assets/appicon.png";
  import {
    makeTwsnmpfcParams,
    makeTwpcapParams,
    makeTwWifiScanParams,
    makeTwWinLogParams,
    makeTwLogEyeParams,
    getConfFromParams,
  } from "./conf.js";
  import TWSNMPFC from "./TWSNMPFC.svelte";
  import TWPCAP from "./TWPCAP.svelte";
  import TWWifiScan from "./TWWifiScan.svelte";
  import TWWinLog from "./TWWinLog.svelte";
  import TWLogEye from "./TWLogEye.svelte";
  import URL from "./URL.svelte";
  import {
    GetInfo,
    GetProcessInfoList,
    Start,
    Stop,
    Delete,
    Save,
    ResetPassword,
  } from "../wailsjs/go/main/App";
  import { BrowserOpenURL } from "../wailsjs/runtime/runtime";
  import {
    Navbar,
    NavBrand,
    Footer,
    FooterCopyright,
    DarkMode,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Table,
    Button,
    Indicator,
    SpeedDial,
    SpeedDialTrigger,
    SpeedDialButton,
    Alert,
    Modal,
    Badge,
  } from "flowbite-svelte";

  let infoMsg = $state("");
  let errorMsg = $state("");
  let oldName = $state("");
  let twsnmpfcModal = $state(false);
  let twsnmpfcConf = $state({});
  let twpcapModal = $state(false);
  let twpcapConf = $state({});
  let twWifiScanModal = $state(false);
  let twWifiScanConf = $state({});
  let twWinLogModal = $state(false);
  let twWinLogConf = $state({});
  let twLogEyeModal = $state(false);
  let twLogEyeConf = $state({});
  let urlModal = $state(false);
  let remoteTwsnmpfcUrl = $state("");

  let info = $state({
    Version: "",
    Env: "",
    NeedPriv: false,
    Ifaces: [],
    PcapVersion: "",
  });

  let proceses = $state([]);
  let waitModal = $state(false);

  onMount(async () => {
    const r = await GetInfo();
    if (r) {
      info = r;
    }
    updateProcessList();
  });

  const updateProcessList = async () => {
    const r = await GetProcessInfoList();
    proceses = r || [];
    proceses.sort((a, b) => {
      if (a.Name < b.Name) {
        return -1;
      }
      if (a.Name > b.Name) {
        return 1;
      }
      return 0;
    });
  };

  const startProcess = async (name, params, task) => {
    errorMsg = "";
    infoMsg = name + "を起動しています。お待ち下さい。";
    waitModal = true;
    if (oldName) {
      await Delete(oldName);
      oldName = "";
    }
    console.log(name, params, task);
    const r = await Start(name, params, task);
    infoMsg = "";
    if (r == "") {
      waitModal = false;
    } else {
      errorMsg = name + "起動エラー:" + r;
    }
    await updateProcessList();
  };

  const stopProcess = async (name) => {
    errorMsg = "";
    infoMsg = name + "を停止しています。お待ち下さい。";
    waitModal = true;
    const r = await Stop(name);
    infoMsg = "";
    if (r == "") {
      waitModal = false;
    } else {
      errorMsg = name + "停止エラー:" + r;
    }
    await updateProcessList();
  };

  const deleteProcess = async (name) => {
    errorMsg = "";
    infoMsg = name + "を削除しています。お待ち下さい。";
    waitModal = true;
    const r = await Delete(name);
    infoMsg = "";
    if (r === "") {
      waitModal = false;
    } else {
      errorMsg = name + "削除エラー:" + r;
    }
    await updateProcessList();
  };

  const resetPassword = async (name) => {
    errorMsg = "";
    infoMsg = name + "のパスワードをリセットしています。お待ち下さい。";
    waitModal = true;
    const r = await ResetPassword(name);
    infoMsg = "";
    if (r === "") {
      waitModal = false;
    } else {
      errorMsg = name + "のパスワードリセットエラー:" + r;
    }
    await updateProcessList();
  };

  const showModal = (name, params, task) => {
    if (name.startsWith("twsnmpfc")) {
      showTwsnmpfcModal(name, params, task);
    } else if (name.startsWith("twpcap")) {
      showTwpcapModal(name, params, task);
    } else if (name.startsWith("twWifiScan")) {
      showTwWifiScanModal(name, params, task);
    } else if (name.startsWith("twWinLog")) {
      showTwWinLogModal(name, params, task);
    } else if (name.startsWith("twlogeye")) {
      showTwLogEyeModal(name, params, task);
    }
  };

  const help = () => {
    BrowserOpenURL("https://zenn.dev/twsnmp/books/twsnmpfc-manual");
  };

  const handleDone = (e) => {
    const detail = e.detail || e;
    switch (detail.type) {
      case "twsnmpfc":
        twsnmpfcModal = false;
        if (detail.start) {
          const name = "twsnmpfc:" + detail.conf.Port;
          startProcess(
            name,
            makeTwsnmpfcParams(detail.conf),
            detail.conf.Task,
          );
        }
        break;
      case "twpcap":
        twpcapModal = false;
        if (detail.start) {
          const name = "twpcap:" + detail.conf.Iface;
          startProcess(
            name,
            makeTwpcapParams(detail.conf),
            detail.conf.Task,
          );
        }
        break;
      case "twWifiScan":
        twWifiScanModal = false;
        if (detail.start) {
          const name = "twWifiScan";
          startProcess(
            name,
            makeTwWifiScanParams(detail.conf),
            detail.conf.Task,
          );
        }
        break;
      case "twWinLog":
        twWinLogModal = false;
        if (detail.start) {
          const name = "twWinLog:" + detail.conf.Remote;
          startProcess(
            name,
            makeTwWinLogParams(detail.conf),
            detail.conf.Task,
          );
        }
        break;
      case "twlogeye":
        twLogEyeModal = false;
        if (detail.start) {
          const name = "twlogeye:" + detail.conf.APIPort;
          startProcess(
            name,
            makeTwLogEyeParams(detail.conf),
            detail.conf.Task,
          );
        }
        break;
      case "url":
        urlModal = false;
        if (detail.url) {
          saveUrl(detail.url);
        }
        break;
    }
  };

  // URLを保存する
  const saveUrl = async (url) => {
    try {
      console.log("saveUrl target:", url);
      if (oldName != "") {
        console.log("deleting oldName:", oldName);
        await Delete(oldName);
      }
      const r = await Save(url);
      console.log("Save result:", r);
      await updateProcessList();
    } catch (err) {
      console.error("saveUrl error:", err);
      errorMsg = "URL保存エラー: " + err;
      waitModal = false;
    }
  };

  const showTwsnmpfcModal = (name, params, task) => {
    if (name == "") {
      twsnmpfcConf = {
        DataStore: "",
        Password: "",
        Local: false,
        Task: false,
        Port: 8080,
      };
      oldName = "";
    } else {
      oldName = name;
      twsnmpfcConf = getConfFromParams(info, params, task);
    }
    twsnmpfcModal = true;
  };

  const showTwpcapModal = (name, params, task) => {
    if (name == "") {
      twpcapConf = {
        Syslog: "",
        Interval: 600,
        Retention: 3600,
        Iface: info.Ifaces[0] ? info.Ifaces[0].Value : "",
        Task: false,
      };
      oldName = "";
    } else {
      oldName = name;
      twpcapConf = getConfFromParams(info, params, task);
    }
    twpcapModal = true;
  };

  const showTwWifiScanModal = (name, params, task) => {
    if (name == "") {
      twWifiScanConf = {
        Syslog: "",
        Interval: 600,
        Task: false,
      };
      oldName = "";
    } else {
      oldName = name;
      twWifiScanConf = getConfFromParams(info, params, task);
    }
    twWifiScanModal = true;
  };

  const showTwWinLogModal = (name, params, task) => {
    if (name == "") {
      twWinLogConf = {
        Syslog: "",
        Interval: 600,
        Remote: "",
        User: "",
        Password: "",
        Task: false,
      };
      oldName = "";
    } else {
      oldName = name;
      twWinLogConf = getConfFromParams(info, params, task);
    }
    twWinLogModal = true;
  };

  const showTwLogEyeModal = (name, params, task) => {
    if (name == "") {
      twLogEyeConf = {
        ConfFile: "",
        APIPort: 8081,
        Key: "",
        Cert: "",
        CACert: "",
        Task: false,
      };
      oldName = "";
    } else {
      oldName = name;
      twLogEyeConf = getConfFromParams(info, params, task);
    }
    twLogEyeModal = true;
  };

  const showUrlModal = (url) => {
    remoteTwsnmpfcUrl = url;
    if (url != "") {
      oldName = url;
    }
    urlModal = true;
  };

  // ブラウザーでTWSNMP FCのURLを開く
  const open = (name, params) => {
    if (name.startsWith("http")) {
      BrowserOpenURL(name);
    } else if (name.startsWith("twsnmpfc:")) {
      const c = getConfFromParams(info, params, false);
      const url = c.TLS
        ? `https://localhost:${c.Port}`
        : `http://localhost:${c.Port}`;
      BrowserOpenURL(url);
    }
  };
</script>

<Navbar>
  <NavBrand>
    <img src={logo} class="mr-3 h-6 sm:h-9" alt="Logo" />
    <span
      class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
    >
      TWSNMP FC 起動ツール
    </span>
    <span class="text-sm dark:text-white ml-2 mr-2">
      {info.Version}
    </span>
    {#if info.NeedPriv}
      <Badge large color="red">管理者権限なし</Badge>
    {/if}
  </NavBrand>
  <div class="flex">
    <DarkMode />
    <Button
      class="!p-2 ml-2 text-xl"
      color="teal"
      onclick={() => updateProcessList()}
    >
      <i class="fa-solid fa-rotate"></i>
    </Button>
    <Button class="!p-2 ml-2 text-xl" color="teal" onclick={help}>
      <i class="fa-regular fa-circle-question"></i>
    </Button>
  </div>
</Navbar>

<Table hoverable={true}>
  <TableHead>
    <TableHeadCell>状態</TableHeadCell>
    <TableHeadCell>名前</TableHeadCell>
    <TableHeadCell>起動/停止</TableHeadCell>
    <TableHeadCell>表示</TableHeadCell>
    <TableHeadCell>PW初期化</TableHeadCell>
    <TableHeadCell>削除</TableHeadCell>
  </TableHead>
  <TableBody class="divide-y">
    {#each proceses as p}
      <TableBodyRow>
        <TableBodyCell>
          {#if p.Running}
            <span class="flex items-center">
              <Indicator size="lg" color="green" class="mr-1.5" />
              稼働
            </span>
          {:else}
            <span class="flex items-center">
              <Indicator size="lg" color="red" class="mr-1.5" />
              停止
            </span>
          {/if}
        </TableBodyCell>
        <TableBodyCell>{p.Name}</TableBodyCell>
        <TableBodyCell>
          {#if p.Name.startsWith("http")}
            <Button
              class="w-2"
              color="teal"
              onclick={() => {
                showUrlModal(p.Name);
              }}
            >
              <i class="fa-solid fa-pencil"></i>
            </Button>
          {:else if p.Running}
            <Button
              class="w-2"
              color="red"
              onclick={() => {
                stopProcess(p.Name);
              }}
            >
              <i class="fa-solid fa-stop"></i>
            </Button>
          {:else}
            <Button
              class="w-2"
              color="teal"
              onclick={() => {
                showModal(p.Name, p.Params, p.Task);
              }}
            >
              <i class="fa-solid fa-play"></i>
            </Button>
          {/if}
        </TableBodyCell>
        <TableBodyCell>
          {#if p.Running && (p.Name.startsWith("twsnmpfc") || p.Name.startsWith("http"))}
            <Button
              class="w-2"
              color="blue"
              onclick={() => {
                open(p.Name, p.Params);
              }}
            >
              <i class="fa-brands fa-chrome"></i>
            </Button>
          {/if}
        </TableBodyCell>
        <TableBodyCell>
          {#if !p.Running && p.Name.startsWith("twsnmpfc")}
            <Button
              class="w-2"
              color="red"
              onclick={() => {
                resetPassword(p.Name);
              }}
            >
              <i class="fa-solid fa-rotate-left"></i>
            </Button>
          {/if}
        </TableBodyCell>
        <TableBodyCell>
          {#if !p.Running || p.Name.startsWith("http")}
            <Button
              class="w-2"
              color="red"
              onclick={() => {
                deleteProcess(p.Name);
              }}
            >
              <i class="fa-solid fa-trash"></i>
            </Button>
          {/if}
        </TableBodyCell>
      </TableBodyRow>
    {/each}
  </TableBody>
</Table>
<SpeedDialTrigger class="absolute right-6 bottom-20" color="teal" />
<SpeedDial tooltip="none" placement="top">
  <SpeedDialButton
    name="twsnmp"
    onclick={() => showTwsnmpfcModal("", [], false)}
  >
    <i class="fa-solid fa-list-check"></i>
  </SpeedDialButton>
  <SpeedDialButton
    name="twpcap"
    onclick={() => showTwpcapModal("", [], false)}
  >
    <i class="fa-solid fa-network-wired"></i>
  </SpeedDialButton>
  <SpeedDialButton
    name="Wifi"
    onclick={() => showTwWifiScanModal("", [], false)}
  >
    <i class="fa-solid fa-wifi"></i>
  </SpeedDialButton>
  {#if info.Env.startsWith("win")}
    <SpeedDialButton
      name="WinLog"
      onclick={() => showTwWinLogModal("", [], false)}
    >
      <i class="fa-brands fa-windows"></i>
    </SpeedDialButton>
  {/if}
  <SpeedDialButton
    name="LogEye"
    onclick={() => showTwLogEyeModal("", [], false)}
  >
    <i class="fa-solid fa-eye"></i>
  </SpeedDialButton>
  <SpeedDialButton name="URL" onclick={() => showUrlModal("")}>
    <i class="fa-brands fa-chrome"></i>
  </SpeedDialButton>
</SpeedDial>

<Footer class="absolute bottom-0 left-0 z-20 w-full">
  <FooterCopyright by="Masayuki Yamai™" year={2021} />
</Footer>

<TWSNMPFC
  {info}
  bind:conf={twsnmpfcConf}
  bind:show={twsnmpfcModal}
  ondone={handleDone}
/>
<TWPCAP {info} bind:conf={twpcapConf} bind:show={twpcapModal} ondone={handleDone} />
<TWWifiScan
  {info}
  bind:conf={twWifiScanConf}
  bind:show={twWifiScanModal}
  ondone={handleDone}
/>
<TWWinLog
  {info}
  bind:conf={twWinLogConf}
  bind:show={twWinLogModal}
  ondone={handleDone}
/>
<TWLogEye
  {info}
  bind:conf={twLogEyeConf}
  bind:show={twLogEyeModal}
  ondone={handleDone}
/>

<URL bind:url={remoteTwsnmpfcUrl} bind:show={urlModal} ondone={handleDone} />

<Modal bind:open={waitModal} size="xs" permanent>
  {#if infoMsg != ""}
    <Alert>
      <i class="fa-solid fa-circle-info"></i>
      {infoMsg}
    </Alert>
  {/if}
  {#if errorMsg != ""}
    <Alert color="red">
      <i class="fa-solid fa-triangle-exclamation"></i>
      {errorMsg}
    </Alert>
    <Button
      color="alternative"
      onmousedown={() => {
        waitModal = false;
        errorMsg = "";
      }}>閉じる</Button
    >
  {/if}
</Modal>
