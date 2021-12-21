<script>
	import "../node_modules/98.css/dist/98.css";
  import Alert from "./Alert.svelte"
  let info = {
    Syslog: '',
    Interval: 600,
    Retention: 3600,
    Status: '',
    Iface: '',
    Ifaces: [],
  }
  let alert = {
    title: '',
    message: '',
  }
  window.go.main.App.GetTWWifiScan().then((r) => {
		info = r
	})
  const start = () => {
		window.go.main.App.Start('twWifiScan','').then((r) => {
			if (r === '') {
        info.Status = '稼働中'
      } else {
        alert.title = 'TWWifiScan起動失敗'
        alert.message = r
      }
		})
  }
  const stop = () => {
		window.go.main.App.Stop('twWifiScan').then((r) => {
			if (r === '') {
        info.Status = '停止'
      } else {
        alert.title = 'TWWifiScan停止失敗'
        alert.message = r
      }
		})
  }
  const addTask = () => {
		window.go.main.App.AddTask('twWifiScan','').then((r) => {
			if (r === '') {
        info.Task = 'Yes'
      } else {
        alert.title = 'TWWifiScanタスク登録失敗'
        alert.message = r
      }
		})
  }
  const delTask = () => {
		window.go.main.App.DelTask('twWifiScan').then((r) => {
			if (r === '') {
        info.Task = 'No'
      } else {
        alert.title = 'TWWifiScanタスク登録解除失敗'
        alert.message = r
      }
		})
  }
</script>

<Alert prop={alert} />
<fieldset>
  <div class="field-row">TWWifiScan:{ info.Status }</div>
  <div class="field-row">
    <label for="syslog">syslog送信先:</label>
    <input id="syslog" type="text" style="width: 80%;" bind:value={info.Syslog} />
  </div>
  <div class="field-row">
    <label for="iface">LAN I/F:</label>
    <select id="iface" bind:value={info.Iface}>
      {#each info.Ifaces as i}
        <option value={i.Value }>{i.Text}</option>
      {/each}
    </select>
  </div>
  <div class="field-row">
    <label for="interval">送信間隔:</label>
    <input id="interval" type="number" min="60" max="3600" bind:value={info.Interval} />
    <label for="interval">秒</label>
  </div>
  <div class="field-row">
    {#if info.Status === '稼働中'}
      <button on:click={stop}>停止</button>
    {:else}
      <button on:click={start}>起動</button>
    {/if}
    {#if info.Task == 'No' }
      <button on:click={addTask}>タスク登録</button>
    {:else if info.Task == 'Yes' }
      <button on:click={delTask}>タスク登録解除</button>
    {/if}
  </div>
</fieldset>

<style>
	label {
		width:  80px;
		margin-left: 10px;
	}
	input[type="number"] {
		width: 50px;
	}
</style>
