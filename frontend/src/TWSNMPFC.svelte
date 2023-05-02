<script>
  import { createEventDispatcher } from "svelte";
  import {
    Button,
    Modal,
    Label,
    Input,
    Toggle,
    ButtonGroup,
    InputAddon,
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
    DataStore: "",
    Password: "",
    Port: 8080,
    Local: false,
    TLS: false,
    Task: false,
  };

  let dataStoreError = false;

  const dispatch = createEventDispatcher();

  const cancel = () => {
    dataStoreError = false;
    dispatch("done", {
      type: "twsnmpfc",
      start: false,
    });
  };

  const start = () => {
    if (!checkParams()) {
      return;
    }
    dispatch("done", {
      type: "twsnmpfc",
      start: true,
      conf: conf,
    });
  };

  const getDataStore = () => {
    window.go.main.App.GetDataStore().then((r) => {
      conf.DataStore = r;
    });
  };

  const checkParams = () => {
    dataStoreError = false;
    if (!conf.DataStore) {
      dataStoreError = true;
      return false;
    }
    return true;
  };
</script>

<Modal bind:open={show} size="md" autoclose={false} class="w-full">
  <h3 class="mb-4 text-xl font-medium text-gray-900 dark:text-white">
    TWSNMP FC起動
  </h3>
  <Label>
    <span>データストア</span>
    <ButtonGroup class="w-full">
      <Input color={dataStoreError ? "red" : "base"} bind:value={conf.DataStore}/>
      <InputAddon>
        <Button class="w-2" on:click={getDataStore}>
          <i class="fa-solid fa-folder" />
        </Button>
      </InputAddon>
    </ButtonGroup>
    <Helper class="mt-2" color={dataStoreError ? "red" : "gray"}>
      TWSNMP FCのデータを保存するディレクトリを指定してください。
    </Helper>
  </Label>
  <Label>
    <span>パスワード</span>
    <Input type="password" bind:value={conf.Password} />
  </Label>
  <div class="flex">
    <Label>
      <span>ポート番号</span>
      <Input type="number" class="w-3" required bind:value={conf.Port} />
    </Label>
    <Toggle class="ml-3 mt-3" bind:checked={conf.Local}>ローカル</Toggle>
    {#if !conf.Local}
      <Toggle class="ml-3 mt-3" bind:checked={conf.TLS}>TLS(HTTPS)</Toggle>
    {/if}
    {#if info.Env == "windows"}
      <Toggle class="ml-3 mt-3" bind:checked={conf.Task}>スケジューラー</Toggle>
    {/if}
  </div>
  <Button on:click={start}>起動</Button>
  <Button color="alternative" on:click={cancel}>キャンセル</Button>
</Modal>
