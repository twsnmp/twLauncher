<script>
	import "../node_modules/98.css/dist/98.css";
	let datastore = '';
	let status = '稼働中'
	let password = "TWSNMPFC!";
	function getDataStore() {
		window.go.main.App.GetDataStore().then((result) => {
			datastore = result
		});
	}
</script>

<fieldset>
  <div class="field-row">TWSNMP FC: {status}</div>
  <div class="field-row">
    <label for="datastore">データストア:</label>
    <input id="datastore" type="text" bind:value={datastore} />
    <button style="min-width: 20px" on:click={getDataStore} >&lt;</button>
    <label for="password">パスワード :</label>
    <input id="password" type="password" style="margin-right: 10px;" bind:value={password} />
    <input type="checkbox" id="local" />
    <label for="local">ローカルモード</label>
  </div>
  <div class="field-row">
    <button>起動</button>
    <button>タスク登録</button>
  </div>
</fieldset>

<style>
	label {
		width:  80px;
		margin-left: 10px;
	}
</style>
