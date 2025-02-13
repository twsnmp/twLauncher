<script>
  import { createEventDispatcher } from "svelte";
  import { Button, Modal, Label, Input, Helper } from "flowbite-svelte";

  export let show = false;
  export let url = "";

  let urlError = false;

  const dispatch = createEventDispatcher();

  const cancel = () => {
    urlError = false;
    dispatch("done", {
      type: "url",
      url: "",
    });
  };

  const save = () => {
    urlError = false;
    const pattern = /^https?:\/\/[\w/:%#\$&\?\(\)~\.=\+\-]+$/
    if (!pattern.test(url)) {
      urlError = true;
      return;
    }
    dispatch("done", {
      type: "url",
      url: url,
    });
  };
</script>

<Modal bind:open={show} size="md" autoclose={false} permanent class="w-full">
  <h3 class="mb-4 text-xl text-gray-900 dark:text-white">リモートTWSNMP FCの設定</h3>
  <Label>
    <span>URL</span>
    <Input
      color={urlError ? "red" : "base"}
      bind:value={url}
      class="mt-1 h-6"
    />
    <Helper class="mt-1" color={urlError ? "red" : "gray"}>
      リモートのTWSNMP FCに接続するためのURLを指定してください。
    </Helper>
  </Label>
  <Button on:click={save}>保存</Button>
  <Button color="alternative" on:click={cancel}>キャンセル</Button>
</Modal>
