<script>
  import { onMount } from 'svelte';
  import { GetSpecs, SaveConfig, CheckCredentials } from '../wailsjs/go/main/App';

  // Use LOWERCASE keys here to match the Go JSON tags
  let specs = {
    cpu: "Loading...",
    ramTotal: "",
    ramModules: "",
    disks: "",
    serial: "",
    tag1: "",
    tag2: "",
    tag3: ""
  };
  
  let statusMessage = "";
  let hasCreds = true;

  async function loadData() {
    hasCreds = await CheckCredentials();
    try {
      const data = await GetSpecs();
      if (data) {
        specs = data; 
      }
    } catch (err) {
      console.error("Error loading specs:", err);
    }
  }

  async function handleSave() {
    statusMessage = "Saving...";
    // Match lowercase keys
    const result = await SaveConfig(specs.tag1, specs.tag2, specs.tag3);
    statusMessage = result;
    setTimeout(() => statusMessage = "", 3000);
  }

  onMount(loadData);
</script>

<main>
  <div class="container">
    {#if !hasCreds}
      <div class="alert-box">
        ⚠️ <strong>Warning:</strong> "service-account.json" not found! 
        Cloud upload will fail. Please place the file in the application folder.
      </div>
    {/if}

    <div class="container">
      <h2>System Information</h2>
    </div>
    
    <div class="card">
      <div class="info-row">
        <span class="label">Serial:</span>
        <span class="value">{specs.serial || "Searching..."}</span>
      </div>
      <div class="info-row">
        <span class="label">CPU:</span>
        <span class="value">{specs.cpu}</span>
      </div>
      <div class="info-row">
        <span class="label">Memory:</span>
        <span class="value">{specs.ramTotal}</span>
      </div>
      <div class="info-row">
        <span class="label">Storage:</span>
        <span class="value">{specs.disks}</span>
      </div>
    </div>

    <h2>Asset Tags</h2>
    <div class="card">
      <div class="input-group">
        <span class="label">TAG_1</span>
        <input bind:value={specs.tag1} />
      </div>
      <div class="input-group">
        <span class="label">TAG_2</span>
        <input bind:value={specs.tag2} />
      </div>
      <div class="input-group">
        <span class="label">TAG_3</span>
        <input bind:value={specs.tag3} />
      </div>
      <button on:click={handleSave}>Save Settings</button>
      {#if statusMessage}<p class="status">{statusMessage}</p>{/if}
    </div>
  </div>
</main>

<style>
  /* FIX: Force light mode theme to prevent "white-on-white" text 
     This overrides Windows Dark Mode defaults in the Wails window.
  */
  :global(body) {
    background-color: #ffffff !important;
    color: #000000 !important;
    margin: 0;
    font-family: sans-serif;
  }

  .container { padding: 20px; max-width: 600px; margin: 0 auto; }
  
  h2 { color: #333; border-bottom: 2px solid #ddd; margin-top: 20px; padding-bottom: 5px; }

  .card { 
    background: #fcfcfc; 
    padding: 15px; 
    border-radius: 8px; 
    border: 1px solid #ccc;
    margin-bottom: 15px;
  }

  .info-row { display: flex; margin-bottom: 8px; border-bottom: 1px solid #eee; padding-bottom: 5px; }
  .label { font-weight: bold; width: 120px; color: #444; }
  .value { color: #000; flex: 1; }

  .input-group { margin-bottom: 12px; display: flex; gap: 10px; align-items: center;}
  
  .alert-box {
    background-color: #fee2e2;
    border: 1px solid #ef4444;
    color: #b91c1c;
    padding: 15px;
    border-radius: 8px;
    margin-bottom: 20px;
    text-align: center;
    font-size: 14px;
  }

  input { 
    width: 100%; 
    padding: 10px; 
    border: 1px solid #bbb; 
    border-radius: 4px; 
    color: #000; 
    background: #fff; 
    box-sizing: border-box;
  }

  button { 
    width: 100%; 
    padding: 12px; 
    background: #2563eb; 
    color: white; 
    border: none; 
    border-radius: 4px; 
    cursor: pointer; 
    font-weight: bold;
    margin-top: 10px;
  }

  button:hover { background: #1d4ed8; }
  
  .status { text-align: center; color: #059669; font-weight: bold; margin-top: 10px; }
</style>