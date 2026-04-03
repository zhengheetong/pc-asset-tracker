<script>
  export let selectedPc = null;
  export let onClose = () => {};
</script>

{#if selectedPc}
  <div class="detail-panel">
    <div class="detail-header">
      <h2>{selectedPc['Serial Number']}</h2>
      <button class="close-btn" on:click={onClose}>✕</button>
    </div>

    <div class="detail-content">
      <div class="detail-group">
        <span class="detail-label">Location / Dept</span>
        <p>{selectedPc['Tag 2']} — {selectedPc['Tag 1']}</p>
      </div>

      <div class="detail-group">
        <span class="detail-label">Operating System</span>
        <p>{selectedPc['Operating System']}</p>
      </div>

      <div class="detail-group">
        <span class="detail-label">Processor (CPU)</span>
        <p>{selectedPc['CPU']}</p>
      </div>

      <div class="detail-group">
        <span class="detail-label">Memory ({selectedPc['Total RAM']})</span>
        <ul class="detail-list">
          {#each (selectedPc['RAM Modules'] || "").split('|') as module}
            <li>{module.trim()}</li>
          {/each}
        </ul>
      </div>

      <div class="detail-group">
        <span class="detail-label">Storage Drives</span>
        <ul class="detail-list">
          {#each (selectedPc['Disks'] || "").split('|') as disk}
            <li>{disk.trim()}</li>
          {/each}
        </ul>
      </div>

      <div class="detail-group">
        <span class="detail-label">Last Audit</span>
        <p class="font-mono">
          {new Date(selectedPc['Timestamp']).toLocaleString()}
        </p>
      </div>
    </div>
  </div>
{/if}
