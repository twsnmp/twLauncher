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
  import {
    GetInfo,
    GetProcessInfoList,
    Start,
    Stop,
    Delete,
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

  let info = {
    Version: "",
    Env: "",
    Ifaces: [],
    PcapVersion: "",
  };

  let proceses = [];

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
  };

  const startProcess = async (name, params, task) => {
    infoMsg = name + "を起動しています。お待ち下さい。";
    if (oldName) {
      await Delete(oldName);
      oldName = "";
    }
    const r = await Start(name, params, task);
    if (r != "") {
      errorMsg = name + "起動エラー:" + r;
    }
    await updateProcessList();
    infoMsg = "";
  };

  const stopProcess = async (name) => {
    infoMsg = name + "を停止しています。お待ち下さい。";
    const r = await Stop(name);
    if (r != "") {
      errorMsg = name + "停止エラー:" + r;
    }
    await updateProcessList();
    infoMsg = "";
  };

  const deleteProcess = async (name) => {
    infoMsg = name + "を削除しています。お待ち下さい。";
    const r = await Delete(name);
    if (r === "") {
      infoMsg = "";
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
    <Button class="!p-2 ml-2 text-xl" color="teal" on:click={help}>
      <i class="fa-regular fa-circle-question" />
    </Button>
  </div>
</Navbar>

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
{/if}

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
          {#if !p.Running}
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
  <SpeedDialButton
    name="WinLog"
    on:click={() => showTwWinLogModal("", [], false)}
  >
    <i class="fa-brands fa-windows" />
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
