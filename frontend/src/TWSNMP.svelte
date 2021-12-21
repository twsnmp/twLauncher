<script>
	import "../node_modules/98.css/dist/98.css"
  import Alert from "./Alert.svelte"
	let info = {
    DataStore:  '',
    Password: '',
    Local: false,
    Status: '',
    Task: '',
  }
  let alert = {
    title: '',
    message: '',
  }
  window.go.main.App.GetTWSNMP().then((r) => {
		info = r
	})
	const getDataStore = () => {
		window.go.main.App.GetDataStore().then((r) => {
			info.DataStore = r
		})
	}
  const start = () => {
		window.go.main.App.Start('twsnmpfc','').then((r) => {
      if (r === '') {
        info.Status = '稼働中'
      } else {
        alert.title = 'TWSNMP起動失敗'
        alert.message = r
      }
    })
  }
  const stop = () => {
		window.go.main.App.Stop('twsnmpfc').then((r) => {
      if (r === '') {
        info.Status = '停止'
      } else {
        alert.title = 'TWSNMP停止失敗'
        alert.message = r
      }
    })
  }
  const addTask = () => {
    window.go.main.App.AddTask('twsnmpfc','').then((r) => {
      if (r==='') {
        info.Task = 'Yes'
      } else {
        alert.title = 'TWSNMPタスク登録失敗'
        alert.message = r
      }
    })
  }
  const delTask = () => {
    window.go.main.App.DelTask('twsnmpfc').then((r) => {
      if (r==='') {
        info.Task = 'No'
      } else {
        alert.title = 'TWSNMPタスク登録解除失敗'
        alert.message = r
      }
    })
  }
</script>

<Alert prop={alert} />
<fieldset>
  <div class="field-row">TWSNMP FC: {info.Status}</div>
  <div class="field-row">
    <label for="datastore">データストア:</label>
    <input id="datastore" type="text" bind:value={info.DataStore} />
    <button style="min-width: 20px" on:click={getDataStore} >&lt;</button>
    <label for="password">パスワード :</label>
    <input id="password" type="password" style="margin-right: 10px;" bind:value={info.Password} />
    <input type="checkbox" id="local" bind:checked={info.Local} />
    <label for="local">ローカルモード</label>
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
</style>
