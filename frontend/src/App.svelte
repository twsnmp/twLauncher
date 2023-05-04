<script>
  import { onMount } from "svelte";
  import logo from "./assets/appicon.png";
  import {
    makeTwsnmpfcParams,
    makeTwpcapParams,
    makeTwWifiScanParams,
    makeTwWinLogParams,
    getConfFromParams,
  } from "./conf.js";
  import TWSNMPFC from "./TWSNMPFC.svelte";
  import TWPCAP from "./TWPCAP.svelte";
  import TWWifiScan from "./TWWifiScan.svelte";
  import TWWinLog from "./TWWinLog.svelte";
  import URL from "./URL.svelte";
  import {
    GetInfo,
    GetProcessInfoList,
    Start,
    Stop,
    Delete,
    Save,
  } from "../wailsjs/go/main/App";
  import {BrowserOpenURL} from "../wailsjs/runtime/runtime";
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
    SpeedDialButton,
    Alert,
    Modal,
  } from "flowbite-svelte";

  let infoMsg = "";
  let errorMsg = "";
  let oldName = "";
  let twsnmpfcModal = false;
  let twsnmpfcConf = {};
  let twpcapModal = false;
  let twpcapConf = {};
  let twWifiScanModal = false;
  let twWifiScanConf = {};
  let twWinLogModal = false;
  let twWinLogConf = {};
  let urlModal = false;
  let remoteTwsnmpfcUrl = "";

  let info = {
    Version: "",
    Env: "",
    Ifaces: [],
    PcapVersion: "",
  };

  let proceses = [];
  let waitModal = false;

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
    proceses.sort((a,b)=> {
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

  const showModal = (name, params, task) => {
    if (name.startsWith("twsnmpfc")) {
      showTwsnmpfcModal(name, params, task);
    } else if (name.startsWith("twpcap")) {
      showTwpcapModal(name, params, task);
    } else if (name.startsWith("twWifiScan")) {
      showTwWifiScanModal(name, params, task);
    } else if (name.startsWith("twWinLog")) {
      showTwWinLogModal(name, params, task);
    }
  };

  const help = () => {
    BrowserOpenURL(
      "https://zenn.dev/twsnmp/books/twsnmpfc-manual"
    );
  };

  const handleDone = (e) => {
    switch (e.detail.type) {
      case "twsnmpfc":
        twsnmpfcModal = false;
        if (e.detail.start) {
          const name = "twsnmpfc:" + e.detail.conf.Port;
          startProcess(
            name,
            makeTwsnmpfcParams(e.detail.conf),
            e.detail.conf.Task
          );
        }
        break;
      case "twpcap":
        twpcapModal = false;
        if (e.detail.start) {
          const name = "twpcap:" + e.detail.conf.Iface;
          startProcess(
            name,
            makeTwpcapParams(e.detail.conf),
            e.detail.conf.Task
          );
        }
        break;
      case "twWifiScan":
        twWifiScanModal = false;
        if (e.detail.start) {
          const name = "twWifiScan";
          startProcess(
            name,
            makeTwWifiScanParams(e.detail.conf),
            e.detail.conf.Task
          );
        }
        break;
      case "twWinLog":
        twWinLogModal = false;
        if (e.detail.start) {
          const name = "twWinLog:" + e.detail.conf.Remote;
          startProcess(
            name,
            makeTwWinLogParams(e.detail.conf),
            e.detail.conf.Task
          );
        }
        break;
      case "url":
        urlModal = false;
        if (e.detail.url) {
          saveUrl(e.detail.url)
        }
        break;
    }
  };

  // URLを保存する
  const saveUrl = async (url)  => {
    if (oldName != "") {
      await Delete(oldName);
    }
    await Save(url);
    updateProcessList();
  }

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

  const showUrlModal = (url) => {
    remoteTwsnmpfcUrl = url;
    if (url != "") {
      oldName = url;
    }
    urlModal = true;
  }

  // ブラウザーでTWSNMP FCのURLを開く
  const open = (name,params) => {
    if (name.startsWith("http")) {
      BrowserOpenURL(name);
    } else if (name.startsWith("twsnmpfc:")) {
      const c = getConfFromParams(info,params,false);
      const url = c.TLS ? `https://localhost:${c.Port}`: `http://localhost:${c.Port}`;
      BrowserOpenURL(url); 
    }
  }
</script>

<Navbar>
  <NavBrand>
    <img src={logo} class="mr-3 h-6 sm:h-9" alt="Logo" />
    <span
      class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
    >
      TWSNMP FC 起動ツール
    </span>
    <span class="text-sm dark:text-white ml-2">
      {info.Version}
    </span>
  </NavBrand>
  <div class="flex">
    <DarkMode />
    <Button class="!p-2 ml-2 text-xl" color="teal" on:click={()=>updateProcessList()}>
      <i class="fa-solid fa-rotate"></i>
    </Button>
    <Button class="!p-2 ml-2 text-xl" color="teal" on:click={help}>
      <i class="fa-regular fa-circle-question" />
    </Button>
  </div>
</Navbar>

<Table hoverable={true}>
  <TableHead>
    <TableHeadCell>状態</TableHeadCell>
    <TableHeadCell>名前</TableHeadCell>
    <TableHeadCell>起動/停止</TableHeadCell>
    <TableHeadCell>表示</TableHeadCell>
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
          {#if p.Name.startsWith("http") }
            <Button
            class="w-2"
            on:click={() => {
              showUrlModal(p.Name);
            }}
          >
            <i class="fa-solid fa-pencil"></i>
          </Button>
        {:else}
            {#if p.Running}
              <Button
                class="w-2"
                color="red"
                on:click={() => {
                  stopProcess(p.Name);
                }}
              >
                <i class="fa-solid fa-stop" />
              </Button>
            {:else}
              <Button
                class="w-2"
                on:click={() => {
                  showModal(p.Name, p.Params, p.Task);
                }}
              >
                <i class="fa-solid fa-play" />
              </Button>
            {/if}
          {/if}
        </TableBodyCell>
        <TableBodyCell>
          {#if p.Running && (p.Name.startsWith("twsnmpfc") || p.Name.startsWith("http")) }
            <Button
              class="w-2"
              color="blue"
              on:click={() => {
                open(p.Name,p.Params);
              }}
            >
              <i class="fa-brands fa-chrome"></i>
            </Button>
          {/if}
        </TableBodyCell>
        <TableBodyCell>
          {#if !p.Running || p.Name.startsWith("http")}
            <Button
              class="w-2"
              color="red"
              on:click={() => {
                deleteProcess(p.Name);
              }}
            >
              <i class="fa-solid fa-trash" />
            </Button>
          {/if}
        </TableBodyCell>
      </TableBodyRow>
    {/each}
  </TableBody>
</Table>
<SpeedDial defaultClass="absolute right-6 bottom-20" tooltip="none">
  <SpeedDialButton
    name="twsnmp"
    on:click={() => showTwsnmpfcModal("", [], false)}
  >
    <i class="fa-solid fa-list-check" />
  </SpeedDialButton>
  <SpeedDialButton
    name="twpcap"
    on:click={() => showTwpcapModal("", [], false)}
  >
    <i class="fa-solid fa-network-wired" />
  </SpeedDialButton>
  <SpeedDialButton 
    name="Wifi"
    on:click={() => showTwWifiScanModal("", [], false)}
  >
    <i class="fa-solid fa-wifi" />
  </SpeedDialButton>
  {#if info.Env.startsWith("win")}
    <SpeedDialButton
      name="WinLog"
      on:click={() => showTwWinLogModal("", [], false)}
    >
      <i class="fa-brands fa-windows" />
    </SpeedDialButton>
  {/if}
  <SpeedDialButton
    name="URL"
    on:click={() => showUrlModal("")}
  >
    <i class="fa-brands fa-chrome" />
  </SpeedDialButton>
</SpeedDial>

<Footer class="absolute bottom-0 left-0 z-20 w-full">
  <FooterCopyright by="Masayuki Yamai™" year={2021} />
</Footer>

<TWSNMPFC
  {info}
  conf={twsnmpfcConf}
  show={twsnmpfcModal}
  on:done={handleDone}
/>
<TWPCAP {info} conf={twpcapConf} show={twpcapModal} on:done={handleDone} />
<TWWifiScan
  {info}
  conf={twWifiScanConf}
  show={twWifiScanModal}
  on:done={handleDone}
/>
<TWWinLog
  {info}
  conf={twWinLogConf}
  show={twWinLogModal}
  on:done={handleDone}
/>

<URL
  url={remoteTwsnmpfcUrl}
  show={urlModal}
  on:done={handleDone}
/>

<Modal bind:open={waitModal} size="xs" autoclose={false} permanent>
  {#if infoMsg != ""}
  <Alert>
    <i class="fa-solid fa-circle-info" />
    {infoMsg}
  </Alert>
  {/if}
  {#if errorMsg != ""}
    <Alert color="red">
      <i class="fa-solid fa-triangle-exclamation" />
      {errorMsg}
    </Alert>
    <Button color="alternative" on:click={()=>{waitModal=false;errorMsg="";}}>閉じる</Button>
  {/if}
</Modal>