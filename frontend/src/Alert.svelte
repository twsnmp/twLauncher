<script>
	import "../node_modules/98.css/dist/98.css"
  export let prop = {
    title: '',
    message: '',
  }
  const close = () => {
    prop.title = ''
    prop.message = ''
  }
</script>

{#if prop.title !== ''}
<div class="modal">
  <div class="window modal-content">
    <div class="title-bar">
      <div class="title-bar-text">{ prop.title }</div>
      <div class="title-bar-controls">
        <button aria-label="Close" on:click={close}></button>
      </div>
    </div>
    <div class="window-body">
      <p>{prop.message}</p>
      <button on:click={close}>OK</button>
    </div>
  </div>
</div>
{/if}

<style>
/* The Modal (background) */
.modal {
  display: block;
  position: fixed;
  z-index: 1;
  left: 0;
  top: 0;
  width: 100%;
  height: 100%;
  overflow: auto;
  background-color: rgb(0,0,0);
  background-color: rgba(0,0,0,0.2);
}
/* Modal Content/Box */
.modal-content {
  margin: 15% auto;
  width: 80%;
}
</style>
