<script>
  import { Button, Modal, Label, Input, Toggle, Helper } from "flowbite-svelte";

  let {
    show = $bindable(false),
    info = {
      Version: "",
      Env: "",
      Ifaces: [],
      PcapVersion: "",
    },
    conf = $bindable({
      Syslog: "",
      Interval: 600,
      Task: false,
    }),
    ondone = undefined,
  } = $props();

  let syslogError = $state(false);

  const cancel = () => {
    syslogError = false;
    show = false;
    if (ondone) {
      ondone({
        type: "twWifiScan",
        start: false,
      });
    }
  };

  const start = () => {
    if (!checkParams()) {
      return;
    }
    show = false;
    if (ondone) {
      ondone({
        type: "twWifiScan",
        start: true,
        conf: conf,
      });
    }
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

<Modal bind:open={show} size="md" class="w-full">
  <h3 class="text-xl text-gray-900 dark:text-white">TWWifiScanの起動</h3>
  <Label>
    <span>syslog送信先</span>
    <Input
      color={syslogError ? "red" : "base"}
      bind:value={conf.Syslog}
      class="mt-1 h-6"
    />
    <Helper class="mt-1" color={syslogError ? "red" : "gray"}>
      syslogの送信先をカンマ区切りで指定してください。
    </Helper>
  </Label>
  <div class="flex">
    <Label>
      <span>送信周期(秒)</span>
      <Input type="number" class="w-3 mt-1 h-6" bind:value={conf.Interval} />
    </Label>
    {#if info.Env == "windows" && !info.NeedPriv}
      <Toggle class="ml-3 mt-4 h-6" bind:checked={conf.Task}>スケジューラー</Toggle>
    {/if}
  </div>
  <Button color="teal" class="mr-2" onmousedown={start}>起動</Button>
  <Button color="alternative" onmousedown={cancel}>キャンセル</Button>
</Modal>
