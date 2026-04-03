<script>
  export let pcs = [];
  export let sortKey = "";
  export let sortAsc = true;
  export let onSort = () => {};
  export let onSelect = () => {};

  function handleSort(key) {
    onSort(key);
  }
</script>

<div class="table-wrapper">
  <table>
    <thead>
      <tr>
        <th on:click={() => handleSort('Serial Number')} class="sortable">
          Serial {sortKey === 'Serial Number' ? (sortAsc ? '↑' : '↓') : ''}
        </th>
        <th on:click={() => handleSort('Operating System')} class="sortable">
          OS {sortKey === 'Operating System' ? (sortAsc ? '↑' : '↓') : ''}
        </th>
        <th on:click={() => handleSort('CPU')} class="sortable">
          CPU {sortKey === 'CPU' ? (sortAsc ? '↑' : '↓') : ''}
        </th>
        <th on:click={() => handleSort('Total RAM')} class="sortable">
          RAM {sortKey === 'Total RAM' ? (sortAsc ? '↑' : '↓') : ''}
        </th>
        <th on:click={() => handleSort('Tag 2')} class="sortable">
          Location {sortKey === 'Tag 2' ? (sortAsc ? '↑' : '↓') : ''}
        </th>
        <th on:click={() => handleSort('Tag 3')} class="sortable">
          Type {sortKey === 'Tag 3' ? (sortAsc ? '↑' : '↓') : ''}
        </th>
      </tr>
    </thead>
    <tbody>
      {#each pcs as pc}
        <tr on:click={() => onSelect(pc)} class="clickable-row">
          <td class="font-mono">{pc['Serial Number']}</td>
          <td>
            <span class="os-badge">
              {pc['Operating System'].replace('Microsoft Windows ', 'Win ')}
            </span>
          </td>
          <td>{pc['CPU']}</td>
          <td class="font-bold">{pc['Total RAM']}</td>
          <td>{pc['Tag 2']}</td>
          <td><span class="type-badge">{pc['Tag 3']}</span></td>
        </tr>
      {/each}
      {#if pcs.length === 0}
        <tr>
          <td colspan="6" class="empty-state">No PCs match your filters.</td>
        </tr>
      {/if}
    </tbody>
  </table>
</div>
