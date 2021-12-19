<script>
	import "../node_modules/98.css/dist/98.css";
	let status = '稼働中'
</script>

<fieldset>
  <div class="field-row">TWPCAP:{ status }</div>
  <div class="field-row">
    <label for="syslog">syslog送信先:</label>
    <input id="syslog" type="text" style="width: 80%;" />
  </div>
  <div class="field-row">
    <label for="iface">LAN I/F:</label>
    <select id="iface">
      <option>eth0</option>
      <option>eth2</option>
      <option>sjsjjsjsjjs</option>
      <option>jsjjsjjsjsjjsjjs</option>
      <option>sjjsjsjjsjsjsjjs</option>
    </select>
  </div>
  <div class="field-row">
    <label for="interval">送信間隔:</label>
    <input id="interval" type="number" min="60" max="3600" value="600" />
    <label for="interval">秒</label>
    <label for="retention">保存時間:</label>
    <input id="retention" type="number" min="600" value="3600" />
    <label for="retention">秒</label>
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
	input[type="number"] {
		width: 50px;
	}
</style>
