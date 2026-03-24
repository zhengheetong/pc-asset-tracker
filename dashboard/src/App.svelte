<script>
  import { onMount } from 'svelte';

  // ⚠️ Remember to paste your actual Google Apps Script URL here!
  const API_URL = "https://script.google.com/macros/s/AKfycbwZZSw_rhSjsP9Q4ogup_L4Yl2-a-afaq_f7ZngTbNxlx2Gt3Ew5IK6sJhizj1WYB6rdQ/exec";
  
  let pcs = [];
  let loading = true;
  let searchQuery = "";
  let selectedLocation = "All";
  let sortKey = "Serial Number"; 
  let sortAsc = true;            
  let currentPage = 1;
  let itemsPerPage = 10;
  let selectedPc = null;

  // 1. DEDUPLICATE: Keep only the newest row for each Serial Number
  $: cleanPcs = extractLatestPcs(pcs);

  // 2. Extract unique locations for the dropdown automatically
  $: locations = ["All", ...new Set(cleanPcs.map(pc => pc['Tag 2']).filter(Boolean))];

  // 3. FILTER: Apply the search box and location dropdown
  $: filteredPcs = cleanPcs.filter(pc => {
    const serial = pc['Serial Number'] || "";
    const cpu = pc['CPU'] || "";
    const matchesSearch = serial.toLowerCase().includes(searchQuery.toLowerCase()) || 
                          cpu.toLowerCase().includes(searchQuery.toLowerCase());
    
    const matchesLocation = selectedLocation === "All" || pc['Tag 2'] === selectedLocation;
    
    return matchesSearch && matchesLocation;
  });

  // 4. SORT: Apply the column sorting to the filtered list
  $: sortedPcs = [...filteredPcs].sort((a, b) => {
    let valA = a[sortKey] || "";
    let valB = b[sortKey] || "";

    // Smart sort for RAM (so "8 GB" comes before "24 GB")
    if (sortKey === 'Total RAM') {
      let numA = parseInt(valA) || 0;
      let numB = parseInt(valB) || 0;
      return sortAsc ? numA - numB : numB - numA;
    }

    // Standard Alphabetical Sort
    if (valA < valB) return sortAsc ? -1 : 1;
    if (valA > valB) return sortAsc ? 1 : -1;
    return 0;
  });

  // --- Functions ---
  // --- NEW: Pagination Logic ---
  // Calculate total pages (minimum 1)
  $: totalPages = Math.max(1, Math.ceil(sortedPcs.length / itemsPerPage));
  
  // UX Fix: If a user searches and the results shrink, force them back to page 1
  $: if (currentPage > totalPages) currentPage = 1;

  // Cut out exactly the 10 rows we need for the current page
  $: paginatedPcs = sortedPcs.slice(
    (currentPage - 1) * itemsPerPage, 
    currentPage * itemsPerPage
  );

  // Pagination Functions
  function nextPage() { if (currentPage < totalPages) currentPage++; }
  function prevPage() { if (currentPage > 1) currentPage--; }

  // Triggered when a user clicks a table header
  function sortBy(key) {
    if (sortKey === key) {
      sortAsc = !sortAsc; // If clicking the same column, reverse the direction
    } else {
      sortKey = key;      // If clicking a new column, switch to it
      sortAsc = true;     // And default to A-Z sorting
    }
  }

  // Filters out the history log, keeping only the most recent scan per PC
  function extractLatestPcs(rawData) {
    const pcMap = new Map();
    rawData.forEach(pc => {
      const serial = pc['Serial Number'];
      if (serial) {
        pcMap.set(serial, pc); 
      }
    });
    return Array.from(pcMap.values());
  }

  // Fetch the data when the page loads
  onMount(async () => {
    try {
      const response = await fetch(API_URL);
      pcs = await response.json();
      loading = false;
    } catch (error) {
      console.error("Failed to load PC data:", error);
      loading = false; // Stop the spinner even if there is an error
    }
  });
</script>

<main class="dashboard-container">
  <header>
    <h1>IT Asset Dashboard</h1>
    <p class="subtitle">Live Hardware Inventory Tracking</p>
  </header>

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Syncing with cloud database...</p>
    </div>
  {:else}
    <div class="card controls-card">
      <div class="input-group">
        <label for="search">Search</label>
        <input id="search" type="text" bind:value={searchQuery} placeholder="Serial or CPU..." />
      </div>

      <div class="input-group">
        <label for="location">Location</label>
        <select id="location" bind:value={selectedLocation}>
          {#each locations as loc}
            <option value={loc}>{loc}</option>
          {/each}
        </select>
      </div>
      
      <div class="stats-badge">
        <span class="pulse"></span>
        {filteredPcs.length} / {cleanPcs.length} Active PCs
      </div>
    </div>

    <div class="card table-card">
      <div class="table-wrapper">
        <table>
          <thead>
            <tr>
              <th on:click={() => sortBy('Serial Number')} class="sortable">
                Serial {sortKey === 'Serial Number' ? (sortAsc ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('Operating System')} class="sortable">
                OS {sortKey === 'Operating System' ? (sortAsc ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('CPU')} class="sortable">
                CPU {sortKey === 'CPU' ? (sortAsc ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('Total RAM')} class="sortable">
                RAM {sortKey === 'Total RAM' ? (sortAsc ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('Tag 2')} class="sortable">
                Location {sortKey === 'Tag 2' ? (sortAsc ? '↑' : '↓') : ''}
              </th>
              <th on:click={() => sortBy('Tag 3')} class="sortable">
                Type {sortKey === 'Tag 3' ? (sortAsc ? '↑' : '↓') : ''}
              </th>
            </tr>
          </thead>
          <tbody>
            {#each paginatedPcs as pc}
              <tr on:click={() => selectedPc = pc} class="clickable-row">
                <td class="font-mono">{pc['Serial Number']}</td>
                <td><span class="os-badge">{pc['Operating System'].replace('Microsoft Windows ', 'Win ')}</span></td>
                <td>{pc['CPU']}</td>
                <td class="font-bold">{pc['Total RAM']}</td>
                <td>{pc['Tag 2']}</td> 
                <td><span class="type-badge">{pc['Tag 3']}</span></td> 
              </tr>
            {/each}
            {#if paginatedPcs.length === 0}
              <tr>
                <td colspan="6" class="empty-state">No PCs match your filters.</td>
              </tr>
            {/if}
            
          </tbody>
        </table>
      </div>
      <div class="pagination-controls">
        <button class="page-btn" on:click={prevPage} disabled={currentPage === 1}>
          &larr; Previous
        </button>
        
        <span class="page-info">
          Page <strong>{currentPage}</strong> of <strong>{totalPages}</strong>
        </span>
        
        <button class="page-btn" on:click={nextPage} disabled={currentPage === totalPages}>
          Next &rarr;
        </button>
      </div>
    </div>
    {#if selectedPc}
      <div class="detail-panel">
        <div class="detail-header">
          <h2>{selectedPc['Serial Number']}</h2>
          <button class="close-btn" on:click={() => selectedPc = null}>✕</button>
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
            <p class="font-mono">{new Date(selectedPc['Timestamp']).toLocaleString()}</p>
          </div>
        </div>
      </div>
    {/if}
  {/if}
</main>

<style>
  /* Global Resets & Dark Mode Variables */
  :global(body) {
    background-color: #0f172a; /* Deep slate background */
    font-family: system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
    color: #f8fafc; /* Off-white text */
    margin: 0;
    padding: 0;
    -webkit-font-smoothing: antialiased;
  }

  /* Colors */
  :root {
    --primary: #3b82f6;
    --primary-hover: #60a5fa;
    --bg-card: #1e293b; /* Slightly lighter slate for cards */
    --border: #334155;  /* Subtle dark border */
    --text-muted: #94a3b8; /* Dimmed text for labels/headers */
  }

  /* Layout */
  .dashboard-container {
    max-width: 100%;
    margin: 2rem 0;
    padding: 0 2rem;
  }

  header { margin-bottom: 2rem; }
  h1 { font-size: 1.8rem; font-weight: 700; margin: 0 0 0.25rem 0; letter-spacing: -0.025em; color: #f8fafc; }
  .subtitle { color: var(--text-muted); margin: 0; font-size: 0.95rem; }

  /* Cards */
  .card {
    background: var(--bg-card);
    border-radius: 12px;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -2px rgba(0, 0, 0, 0.3);
    border: 1px solid var(--border);
    margin-bottom: 1.5rem;
  }

  /* Controls */
  .controls-card {
    display: flex;
    flex-wrap: wrap;
    gap: 1.5rem;
    padding: 1.25rem 1.5rem;
    align-items: flex-end;
  }

  .input-group {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    flex: 1;
    min-width: 200px;
  }

  .input-group label { font-size: 0.85rem; font-weight: 600; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
  
  input, select {
    padding: 0.6rem 1rem;
    border: 1px solid var(--border);
    border-radius: 8px;
    font-size: 0.95rem;
    transition: all 0.2s ease;
    background-color: #0f172a; /* Dark inputs */
    color: #f8fafc;
  }
  
  input:focus, select:focus {
    outline: none;
    border-color: var(--primary);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.25);
    background-color: #1e293b;
  }

  /* Live Stats Badge */
  .stats-badge {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.6rem 1rem;
    background: rgba(59, 130, 246, 0.1); /* Transparent blue */
    color: #60a5fa;
    border-radius: 8px;
    font-weight: 600;
    font-size: 0.9rem;
    margin-left: auto;
    border: 1px solid rgba(59, 130, 246, 0.2);
  }
  .pulse {
    width: 8px;
    height: 8px;
    background-color: #3b82f6;
    border-radius: 50%;
    animation: pulse-animation 2s infinite;
  }
  @keyframes pulse-animation {
    0% { box-shadow: 0 0 0 0 rgba(59, 130, 246, 0.7); }
    70% { box-shadow: 0 0 0 6px rgba(59, 130, 246, 0); }
    100% { box-shadow: 0 0 0 0 rgba(59, 130, 246, 0); }
  }

  /* Table */
  .table-wrapper { overflow-x: auto; }
  table { width: 100%; border-collapse: collapse; text-align: left; }
  
  th {
    padding: 1rem 1.5rem;
    font-size: 0.85rem;
    font-weight: 600;
    color: var(--text-muted);
    border-bottom: 2px solid var(--border);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    background-color: #0f172a; /* Darkest background for headers */
  }
  th.sortable { cursor: pointer; user-select: none; transition: background 0.2s; color: #cbd5e1; }
  th.sortable:hover { background-color: #1e293b; color: #f8fafc; }

  td {
    padding: 1rem 1.5rem;
    border-bottom: 1px solid var(--border);
    font-size: 0.95rem;
    color: #cbd5e1; /* Light gray text for table data */
    white-space: nowrap;
  }
  tr:last-child td { border-bottom: none; }
  tr:hover { background-color: #334155; /* Highlight row on hover */ }

  /* Typography Enhancements */
  .font-mono { font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 0.9rem; color: #94a3b8; }
  .font-bold { font-weight: 600; color: #f8fafc; }

  /* Badges */
  .os-badge { background: #334155; padding: 0.25rem 0.6rem; border-radius: 6px; font-size: 0.85rem; font-weight: 500; border: 1px solid #475569; color: #e2e8f0; }
  .type-badge { background: rgba(168, 85, 247, 0.15); color: #c084fc; padding: 0.25rem 0.6rem; border-radius: 6px; font-size: 0.85rem; font-weight: 600; border: 1px solid rgba(168, 85, 247, 0.3); }

  /* Loading & Empty States */
  .loading-state { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem; color: var(--text-muted); }
  .spinner { border: 3px solid #334155; border-top: 3px solid var(--primary); border-radius: 50%; width: 24px; height: 24px; animation: spin 1s linear infinite; margin-bottom: 1rem; }
  @keyframes spin { 0% { transform: rotate(0deg); } 100% { transform: rotate(360deg); } }
  .empty-state { text-align: center; padding: 3rem; color: var(--text-muted); font-style: italic; }

  /* Pagination Controls */
  .pagination-controls {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-top: 1px solid var(--border);
    background-color: #0f172a;
    border-bottom-left-radius: 12px;
    border-bottom-right-radius: 12px;
  }

  .page-info { font-size: 0.9rem; color: var(--text-muted); }
  .page-info strong { color: #f8fafc; }

  .page-btn {
    padding: 0.5rem 1rem;
    background-color: #1e293b;
    border: 1px solid var(--border);
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    color: #e2e8f0;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .page-btn:hover:not(:disabled) {
    background-color: #334155;
    border-color: #475569;
    color: #f8fafc;
  }

  .page-btn:disabled {
    opacity: 0.4;
    cursor: not-allowed;
    background-color: #0f172a;
  }

  /* Clickable Table Rows */
  .clickable-row {
    cursor: pointer;
    transition: background-color 0.15s ease;
  }
  .clickable-row:hover {
    background-color: #334155; /* Highlight color from the dark theme */
  }

  /* Floating Detail Panel */
  .detail-panel {
    position: fixed;
    top: 2rem;
    right: 2rem;
    width: 380px;
    max-height: calc(100vh - 4rem);
    overflow-y: auto;
    background-color: var(--bg-card);
    border: 1px solid var(--border);
    border-radius: 12px;
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.5), 0 8px 10px -6px rgba(0, 0, 0, 0.5);
    z-index: 100;
    animation: slideIn 0.3s cubic-bezier(0.16, 1, 0.3, 1);
  }

  @keyframes slideIn {
    from { opacity: 0; transform: translateX(20px) scale(0.98); }
    to { opacity: 1; transform: translateX(0) scale(1); }
  }

  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1.25rem 1.5rem;
    border-bottom: 1px solid var(--border);
    background-color: #0f172a;
    position: sticky;
    top: 0;
  }

  .detail-header h2 {
    margin: 0;
    font-size: 1.25rem;
    font-family: ui-monospace, monospace;
    color: var(--primary-hover);
  }

  .close-btn {
    background: none;
    border: none;
    color: var(--text-muted);
    font-size: 1.25rem;
    cursor: pointer;
    padding: 0.25rem;
    line-height: 1;
    transition: color 0.2s;
  }

  .close-btn:hover {
    color: #ef4444; /* Red close button on hover */
  }

  .detail-content {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .detail-label {
    display: block;
    font-size: 0.75rem;
    font-weight: 700;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 0.25rem;
  }

  .detail-group p {
    margin: 0;
    font-size: 0.95rem;
    color: #e2e8f0;
    line-height: 1.4;
  }

  .detail-list {
    margin: 0;
    padding-left: 1.25rem;
    color: #e2e8f0;
    font-size: 0.9rem;
  }
  
  .detail-list li {
    margin-bottom: 0.25rem;
  }
</style>