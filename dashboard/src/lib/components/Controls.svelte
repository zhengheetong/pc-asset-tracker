<script>
  import StatsBadge from './StatsBadge.svelte';

  export let searchQuery = "";
  export let selectedLocation = "All";
  export let locations = [];
  export let filteredCount = 0;
  export let totalCount = 0;
</script>

<div class="card controls-card">
  <div class="input-group">
    <label for="search">Search</label>
    <input
      id="search"
      type="text"
      bind:value={searchQuery}
      placeholder="Serial or CPU..."
    />
  </div>

  <div class="input-group">
    <label for="location">Location</label>
    <select id="location" bind:value={selectedLocation}>
      {#each locations as loc}
        <option value={loc}>{loc}</option>
      {/each}
    </select>
  </div>

  <StatsBadge count={filteredCount} total={totalCount} />
</div>
