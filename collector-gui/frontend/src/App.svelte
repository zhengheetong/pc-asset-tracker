<script>
  import { onMount } from 'svelte';
  import { GetSpecs, SaveConfig, CheckCredentials, InstallToPC } from '../wailsjs/go/main/App';

  let specs = {
    os: "Scanning...",
    cpu: "Loading...",
    ramTotal: "",
    ramModules: "",
    disks: "",
    serial: "",
    tag1: "Default",
    tag2: "Default",
    tag3: "Default"
  };

  let statusMessage = "";
  let hasCreds = false;
  let isInstalling = false;
  let installMessage = "";
  let isError = false;

  async function loadData() {
    hasCreds = await CheckCredentials();
    try {
      const data = await GetSpecs();
      if (data) {
        specs = data;
        // Default tags if none exist
        if (!specs.tag1) specs.tag1 = "Default";
        if (!specs.tag2) specs.tag2 = "Default";
        if (!specs.tag3) specs.tag3 = "Default";
      }
    } catch (err) {
      console.error("Error loading specs:", err);
    }
  }

  async function handleSave() {
    statusMessage = "Saving...";
    const result = await SaveConfig(specs.tag1, specs.tag2, specs.tag3);
    statusMessage = result;
    setTimeout(() => statusMessage = "", 3000);
  }

  async function handleInstall() {
    isInstalling = true;
    installMessage = "";
    isError = false;

    try {
      const result = await InstallToPC();
      installMessage = result; 
    } catch (err) {
      isError = true;
      installMessage = "Installation failed: " + err;
    } finally {
      isInstalling = false;
    }
  }

  onMount(loadData);
</script>

<main data-wails-drag>
  <div class="app-container" data-wails-no-drag>
    
    <header>
      <h1>PC Tracker Agent</h1>
      <p class="subtitle">Background Hardware Collector</p>
    </header>

    {#if !hasCreds}
      <div class="alert-box">
        ⚠️ <strong>Warning:</strong> "service-account.json" not found in folder!
      </div>
    {/if}

    <div class="dashboard-grid">
      
      <div class="left-col">
        <div class="card system-card">
          <h2 class="text-center">System Information</h2>
          
          <div class="info-row">
            <span class="label">Serial:</span>
            <span class="value font-mono">{specs.serial || "Searching..."}</span>
          </div>
          <div class="info-row">
            <span class="label">OS:</span>
            <span class="value">{specs.os}</span>
          </div>
          <div class="info-row">
            <span class="label">CPU:</span>
            <span class="value">{specs.cpu}</span>
          </div>
          <div class="info-row">
            <span class="label">Memory:</span>
            <span class="value">{specs.ramTotal}</span>
          </div>
          
          <div class="info-row align-top">
            <span class="label">RAM:</span>
            <div class="value list-value">
              {#each (specs.ramModules || "").split('|') as module}
                {#if module.trim()}
                  <div class="list-item">{module.trim()}</div>
                {/if}
              {/each}
            </div>
          </div>
          
          <div class="info-row align-top border-none">
            <span class="label">Storage:</span>
            <div class="value list-value">
              {#each (specs.disks || "").split('|') as disk}
                {#if disk.trim()}
                  <div class="list-item">{disk.trim()}</div>
                {/if}
              {/each}
            </div>
          </div>
        </div>
      </div>

      <div class="right-col">
        
        <div class="card">
          <h2 class="text-center">Asset Tags</h2>
          
          <div class="input-group">
            <span class="label">Department</span>
            <input bind:value={specs.tag1} placeholder="Default" />
          </div>
          <div class="input-group">
            <span class="label">Location</span>
            <input bind:value={specs.tag2} placeholder="Default" />
          </div>
          <div class="input-group">
            <span class="label">Type</span>
            <input bind:value={specs.tag3} placeholder="Default" />
          </div>
          
          <button class="primary-btn" on:click={handleSave}>Save Tags to Config</button>
          {#if statusMessage}<p class="status-msg success-text">{statusMessage}</p>{/if}
        </div>

        <div class="card install-card">
          <h2 class="text-center">Deployment</h2>
          
          <button 
            class="install-btn" 
            on:click={handleInstall} 
            disabled={!hasCreds || isInstalling || installMessage.includes("Successfully")}
          >
            {#if isInstalling}
              Installing to Background...
            {:else if installMessage.includes("Successfully")}
              Installed & Active ✓
            {:else}
              Install Silent Tracker
            {/if}
          </button>

          {#if installMessage}
            <p class="status-msg {isError ? 'error-text' : 'success-text'}">
              {installMessage}
            </p>
          {/if}
        </div>
        
      </div>
    </div>
  </div>
</main>

<style>
  /* Dark Mode Slate Theme */
  :global(body) {
    background-color: #0f172a; /* Main background */
    color: #f8fafc;
    margin: 0;
    font-family: system-ui, -apple-system, sans-serif;
    -webkit-user-select: none;
    height: 100vh;
    overflow-y: auto;
  }

  main {
    padding: 2rem;
    display: flex;
    justify-content: center;
    align-items: center;
    min-height: 100vh;
    box-sizing: border-box;
  }

  .app-container {
    width: 100%;
    max-width: 900px; /* Widened to fit two columns */
  }

  header { text-align: center; margin-bottom: 2rem; }
  h1 { margin: 0; font-size: 1.8rem; font-weight: 600; color: #f8fafc; }
  .subtitle { color: #94a3b8; margin: 0.25rem 0 0 0; font-size: 1rem; }

  /* CSS Grid Layout */
  .dashboard-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    align-items: flex-start;
  }

  .left-col, .right-col {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  /* Cards */
  .card { 
    background: #1e293b; /* Card background matching the image */
    padding: 1.5rem;
    border-radius: 8px; 
    border: 1px solid #334155;
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3);
  }

  .system-card { padding-bottom: 0.5rem; } /* Tighter bottom spacing for left card */

  h2.text-center { 
    color: #f8fafc; 
    font-size: 1.1rem; 
    margin: 0 0 1.25rem 0; 
    padding-bottom: 0.75rem; 
    border-bottom: 1px solid #334155; 
    text-align: center;
  }

  /* Info Rows (Left aligned text) */
  .info-row { 
    display: flex; 
    align-items: center;
    margin-bottom: 1rem; 
    border-bottom: 1px solid #273548;
    padding-bottom: 1rem; 
    font-size: 0.9rem;
  }
  .info-row.align-top { align-items: flex-start; }
  .info-row.border-none { border-bottom: none; }
  
  .label { 
    font-weight: 600; 
    width: 90px; 
    color: #94a3b8; 
    flex-shrink: 0; 
    text-align: left; /* Forced Left Alignment */
  }
  
  .value { 
    color: #e2e8f0; 
    flex: 1; 
    word-break: break-word; 
    text-align: left; /* Forced Left Alignment */
  }
  
  .font-mono { font-family: ui-monospace, monospace; color: #60a5fa; }

  /* Lists inside values (RAM and Storage) */
  .list-value {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .list-item { line-height: 1.4; }

  /* Inputs */
  .input-group { 
    margin-bottom: 1rem; 
    display: flex; 
    align-items: center;
    gap: 15px;
  }
  
  input { 
    flex: 1; 
    padding: 0.6rem 0.75rem;
    border: 1px solid #334155; 
    border-radius: 4px; 
    color: #f8fafc; 
    background: #0f172a; 
    font-size: 0.9rem;
    transition: border-color 0.2s;
  }
  input:focus { outline: none; border-color: #3b82f6; }

  /* Buttons */
  button { 
    width: 100%; 
    padding: 0.85rem; 
    border: none; 
    border-radius: 6px;
    cursor: pointer; 
    font-weight: 600;
    font-size: 0.95rem;
    transition: all 0.2s ease;
    margin-top: 0.5rem;
  }

  .primary-btn { 
    background: #334155; 
    color: #f8fafc; 
    border: 1px solid #475569; 
  }
  .primary-btn:hover { background: #475569; }

  .install-card { background: #111827; border-color: #1f2937; }
  .install-btn { background: #3b82f6; color: white; }
  .install-btn:hover:not(:disabled) { background: #2563eb; }
  .install-btn:disabled { background: #1e293b; color: #64748b; cursor: not-allowed; }

  /* Alerts */
  .alert-box {
    background-color: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.4);
    color: #f87171;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    text-align: center;
    font-size: 0.9rem;
  }

  .status-msg { text-align: center; font-weight: 500; margin: 0.75rem 0 0 0; font-size: 0.9rem; }
  .success-text { color: #4ade80; }
  .error-text { color: #f87171; }

  /* Responsive styling for smaller screens */
  @media (max-width: 768px) {
    .dashboard-grid { grid-template-columns: 1fr; }
    .app-container { max-width: 500px; }
  }
</style>