<script>
  import { Button, Modal, Label, Input, Helper } from "flowbite-svelte";

  let {
    show = $bindable(false),
    url = $bindable(""),
    ondone = undefined,
  } = $props();

  let urlError = $state(false);

  const cancel = () => {
    urlError = false;
    show = false;
    if (ondone) {
      ondone({
        type: "url",
        url: "",
      });
    }
  };

  const save = () => {
    urlError = false;
    const pattern = /^https?:\/\/[\w/:%#\$&\?\(\)~\.=\+\-]+$/;
    if (!pattern.test(url)) {
      urlError = true;
      return;
    }
    show = false;
    if (ondone) {
      ondone({
        type: "url",
        url: url,
      });
    }
  };
</script>

<Modal bind:open={show} size="md" class="w-full">
  <h3 class="mb-4 text-xl text-gray-900 dark:text-white">
    リモートTWSNMP FCの設定
  </h3>
  <Label>
    <span>URL</span>
    <Input
      color={urlError ? "red" : "default"}
      bind:value={url}
      class="mt-2"
    />
    <Helper class="mt-2" color={urlError ? "red" : "gray"}>
      リモートのTWSNMP FCに接続するためのURLを指定してください。
    </Helper>
  </Label>
  <Button color="teal" class="mr-2" onmousedown={save}>保存</Button>
  <Button color="alternative" onmousedown={cancel}>キャンセル</Button>
</Modal>
