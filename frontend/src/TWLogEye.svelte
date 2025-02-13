<script>
  import { createEventDispatcher } from "svelte";
  import {
    Button,
    Modal,
    Label,
    Input,
    Toggle,
    Helper,
    ButtonGroup,
    InputAddon,
  } from "flowbite-svelte";
  import { GetFile } from "../wailsjs/go/main/App";
  export let show = false;
  export let info = {
    Version: "",
    Env: "",
    Ifaces: [],
    PcapVersion: "",
    NeedPriv: false,
  };

  export let conf = {
    APIPort: 8081,
    ConfFile: "",
    Cert: "",
    Key: "",
    CACert: "",
    Task: false,
  };

  var portError = false;
  var confFileError = false;

  const dispatch = createEventDispatcher();

  const cancel = () => {
    dispatch("done", {
      type: "twlogeye",
      start: false,
    });
  };

  const start = () => {
    if (!checkParams()) {
      return;
    }
    dispatch("done", {
      type: "twlogeye",
      start: true,
      conf: conf,
    });
  };

  const checkParams = () => {
    confFileError = false;
    portError = false;
    if (!conf.ConfFile) {
      confFileError = true;
      return false;
    }
    if (conf.Port < 1 || conf.Port > 0xfffe) {
      portError = true;
      return false;
    }    
    return true;
  };

  const getConfFile = async () => {
    conf.ConfFile = await GetFile("設定ファイル");
  };

  const getKeyFile = async () => {
    conf.Key = await GetFile("APIサーバー秘密鍵ファイル");
  };

  const getCertFile = async () => {
    conf.Cert = await GetFile("APIサーバー証明書ファイル");
  };

  const getCACertFile = async () => {
    conf.Cert = await GetFile("APIサーバーCA証明書ファイル");
  };
</script>

<Modal bind:open={show} size="md" autoclose={false} permanent class="w-full">
  <h3 class="text-xl text-gray-900 dark:text-white">TWLogEyeの起動</h3>
  <div class="grid grid-cols-3 gap-2">
    <Label class="col-span-2">
      <span>設定ファイル</span>
      <ButtonGroup class="w-full">
        <Input
          class="h-6"
          color={confFileError ? "red" : "base"}
          bind:value={conf.ConfFile}
        />
        <InputAddon>
          <Button class="w-1 h-5" on:click={getConfFile}>
            <i class="fa-solid fa-file" />
          </Button>
        </InputAddon>
      </ButtonGroup>
      <Helper class="mt-1" color={confFileError ? "red" : "gray"}>
        設定ファイルを指定してください。
      </Helper>
    </Label>
    <Label>
      <span>APIポート番号</span>
      <Input
        color={portError ? "red" : "base"}
        type="number"
        class="ml-2 w-2 h-6"
        bind:value={conf.APIPort}
      />
      <Helper class="mt-1" color={portError ? "red" : "gray"}>
        APIサーバーのポート番号を指定してください。
      </Helper>
    </Label>
  </div>
  <Label>
    <span>APIサーバー秘密鍵ファイル</span>
    <ButtonGroup class="w-full">
      <Input class="h-6" bind:value={conf.Key} />
      <InputAddon>
        <Button class="w-1 h-5" on:click={getKeyFile}>
          <i class="fa-solid fa-key" />
        </Button>
      </InputAddon>
    </ButtonGroup>
    <Helper class="mt-1">APIサーバー秘密鍵ファイルを指定してください。</Helper>
  </Label>
  <Label>
    <span>APIサーバー証明書ファイル</span>
    <ButtonGroup class="w-full">
      <Input class="h-6" bind:value={conf.Cert} />
      <InputAddon>
        <Button class="w-2 h-5" on:click={getCertFile}>
          <i class="fa-solid fa-certificate" />
        </Button>
      </InputAddon>
    </ButtonGroup>
    <Helper class="mt-1">APIサーバー証明書ファイルを指定してください。</Helper>
  </Label>
  <Label>
    <span>APIサーバーCA証明書ファイル</span>
    <ButtonGroup class="w-full">
      <Input class="h-6" bind:value={conf.CACert} />
      <InputAddon>
        <Button class="w-2 h-5" on:click={getCACertFile}>
          <i class="fa-solid fa-certificate" />
        </Button>
      </InputAddon>
    </ButtonGroup>
    <Helper class="mt-1">
      APIサーバーのCA証明書ファイルを指定してください。
    </Helper>
  </Label>
  {#if info.Env == "windows" && !info.NeedPriv}
    <Toggle class="mt-2 h-6" bind:checked={conf.Task}>スケジューラー</Toggle>
  {/if}
  <Button on:click={start}>起動</Button>
  <Button color="alternative" on:click={cancel}>キャンセル</Button>
</Modal>
