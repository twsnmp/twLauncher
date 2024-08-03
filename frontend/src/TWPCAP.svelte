<script>
  import { createEventDispatcher } from "svelte";
  import {
    Button,
    Modal,
    Label,
    Input,
    Toggle,
    Alert,
    Select,
    Helper,
  } from "flowbite-svelte";

  export let show = false;
  export let info = {
    Version: "",
    Env: "",
    Ifaces: [],
    PcapVersion: "",
  };

  export let conf = {
    Syslog: "",
    Interval: 600,
    Retention: 3600,
    Iface: "",
    Task: false,
  };

  let syslogError = false;
  let ifaceError = false;

  const dispatch = createEventDispatcher();

  const cancel = () => {
    syslogError = false;
    ifaceError = false;
    dispatch("done", {
      type: "twpcap",
      start: false,
    });
  };

  const start = () => {
    if (!checkParams()) {
      return;
    }
    dispatch("done", {
      type: "twpcap",
      start: true,
      conf: conf,
    });
  };
  const checkParams = () => {
    syslogError = false;
    ifaceError = false;
    if (!conf.Syslog) {
      syslogError = true;
      return false;
    }
    if (!conf.Iface) {
      ifaceError = true;
      return false;
    }
    return true;
  };
  const help = () => {
    window.runtime.BrowserOpenURL(
      "https://zenn.dev/twsnmp/books/twsnmpfc-manual/viewer/sensors"
    );
  };
</script>

<Modal bind:open={show} size="md" autoclose={false} permanent class="w-full">
  <h3 class="text-xl font-medium text-gray-900 dark:text-white">
    twpcapの起動
  </h3>
  {#if info.PcapVersion}
    <Alert color="blue">
      PCAP Version:{info.PcapVersion}
    </Alert>
  {:else}
    <Alert color="red">
      <i class="fa-solid fa-triangle-exclamation" />
      PCAPライブラリーをインストールしてください。
      <Button class="w-2" color="alternative" on:click={help}>
        <i class="fa-regular fa-circle-question" />
      </Button>
    </Alert>
  {/if}
  <Label>
    <span>syslog送信先</span>
    <Input
      color={syslogError ? "red" : "base"}
      bind:value={conf.Syslog}
      class="mt-1"
    />
    <Helper class="mt-2" color={syslogError ? "red" : "gray"}>
      syslogの送信先をカンマ区切りで指定してください。
    </Helper>
  </Label>
  <Label
    >LANポート
    <Select
      class="mt-1"
      color={ifaceError ? "red" : "base"}
      bind:value={conf.Iface}
      placeholder=""
    >
      {#each info.Ifaces as i}
        {#if i.Value == conf.Iface}
          <option selected value={i.Value}>{i.Text}</option>
        {:else}
          <option value={i.Value}>{i.Text}</option>
        {/if}
      {/each}
    </Select>
    <Helper class="mt-2" color={ifaceError ? "red" : "gray"}>
      パケットキャプチャーするLANポートを指定してください。
    </Helper>
  </Label>

  <div class="flex">
    <Label>
      <span>送信周期(秒)</span>
      <Input type="number" class="w-2 mt-1" bind:value={conf.Interval} />
    </Label>
    <Label class="ml-3">
      <span>保存時間(秒)</span>
      <Input type="number" class="w-3 mt-1" bind:value={conf.Retention} />
    </Label>
    {#if info.Env == "windows" && !info.NeedPriv}
      <Toggle class="ml-3 mt-3" bind:checked={conf.Task}>スケジューラー</Toggle>
    {/if}
  </div>
  {#if info.PcapVersion}
    <Button on:click={start}>起動</Button>
  {/if}
  <Button color="alternative" on:click={cancel}>キャンセル</Button>
</Modal>
