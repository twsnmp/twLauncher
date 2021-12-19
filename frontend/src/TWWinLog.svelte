<script>
	import "../node_modules/98.css/dist/98.css";
	let status = '稼働中'
</script>

<fieldset>
  <div class="field-row">TWWinLog:{ status }</div>
  <div class="field-row">
    <label for="twWinLog_syslog">syslog送信先:</label>
    <input id="twWinLog_syslog" type="text" style="width: 80%;" />
  </div>
  <div class="field-row">
    <label for="twWinLog_remote">リモート:</label>
    <input id="twWinLog_remote" type="text"/>
    <label for="twWinLog_user">ユーザー:</label>
    <input id="twWinLog_user" type="text"/>
    <label for="twWinLog_password">パスワード:</label>
    <input id="twWinLog_password" type="password"/>
    <label for="twWinLog_interval">送信間隔:</label>
    <input id="twWinLog_interval" type="number" min="60" max="3600" value="600" />
    <label for="twWinLog_interval">秒</label>
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
