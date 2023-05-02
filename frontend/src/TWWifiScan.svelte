<script>
  import { createEventDispatcher } from "svelte";
  import { Button, Modal, Label, Input, Toggle, Helper } from "flowbite-svelte";

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
    Task: false,
  };

  let syslogError = false;

  const dispatch = createEventDispatcher();

  const cancel = () => {
    syslogError = false;
    dispatch("done", {
      type: "twWifiScan",
      start: false,
    });
  };

  const start = () => {
    if (!checkParams()) {
      return;
    }
    dispatch("done", {
      type: "twWifiScan",
      start: true,
      conf: conf,
    });
  };
  const checkParams = () => {
    syslogError = false;
    if (!conf.Syslog) {
      syslogError = true;
      return false;
    }
    return true;
  };
</script>

<Modal bind:open={show} size="md" autoclose={false} permanent class="w-full">
  <h3 class="mb-4 text-xl text-gray-900 dark:text-white">twWifiScanの設定</h3>
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
  <div class="flex">
    <Label>
      <span>送信周期(秒)</span>
      <Input type="number" class="w-3 mt-1" bind:value={conf.Interval} />
    </Label>
    {#if info.Env == "windows"}
      <Toggle class="ml-3 mt-3" bind:checked={conf.Task}>スケジューラー</Toggle>
    {/if}
  </div>
  <Button on:click={start}>起動</Button>
  <Button color="alternative" on:click={cancel}>キャンセル</Button>
</Modal>
