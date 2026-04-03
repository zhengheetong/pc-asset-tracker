<script>
  import { onMount } from "svelte";
  import Header from "./lib/components/Header.svelte";
  import Controls from "./lib/components/Controls.svelte";
  import DeviceTable from "./lib/components/DeviceTable.svelte";
  import DeviceDetail from "./lib/components/DeviceDetail.svelte";
  import Pagination from "./lib/components/Pagination.svelte";

  // Google Apps Script URL
  const API_URL =
    "https://script.google.com/macros/s/AKfycbwZZSw_rhSjsP9Q4ogup_L4Yl2-a-afaq_f7ZngTbNxlx2Gt3Ew5IK6sJhizj1WYB6rdQ/exec";

  let pcs = [];
  let loading = true;
  let searchQuery = "";
  let selectedLocation = "All";
  let sortKey = "Serial Number";
  let sortAsc = true;
  let currentPage = 1;
  let itemsPerPage = 10;
  let selectedPc = null;

  // --- Reactive Logic ---

  // 1. DEDUPLICATE: Keep only the newest row for each Serial Number
  $: cleanPcs = extractLatestPcs(pcs);

  // 2. Extract unique locations for the dropdown automatically
  $: locations = [
    "All",
    ...new Set(cleanPcs.map((pc) => pc["Tag 2"]).filter(Boolean)),
  ];

  // 3. FILTER: Apply the search box and location dropdown
  $: filteredPcs = cleanPcs.filter((pc) => {
    const serial = pc["Serial Number"] || "";
    const cpu = pc["CPU"] || "";
    const matchesSearch =
      serial.toLowerCase().includes(searchQuery.toLowerCase()) ||
      cpu.toLowerCase().includes(searchQuery.toLowerCase());

    const matchesLocation =
      selectedLocation === "All" || pc["Tag 2"] === selectedLocation;

    return matchesSearch && matchesLocation;
  });

  // 4. SORT: Apply column sorting
  $: sortedPcs = [...filteredPcs].sort((a, b) => {
    let valA = a[sortKey] || "";
    let valB = b[sortKey] || "";

    if (sortKey === "Total RAM") {
      let numA = parseInt(valA) || 0;
      let numB = parseInt(valB) || 0;
      return sortAsc ? numA - numB : numB - numA;
    }

    if (valA < valB) return sortAsc ? -1 : 1;
    if (valA > valB) return sortAsc ? 1 : -1;
    return 0;
  });

  // 5. PAGINATION: Slice data for current page
  $: totalPages = Math.max(1, Math.ceil(sortedPcs.length / itemsPerPage));
  $: if (currentPage > totalPages) currentPage = 1;
  $: paginatedPcs = sortedPcs.slice(
    (currentPage - 1) * itemsPerPage,
    currentPage * itemsPerPage
  );

  // --- Functions ---

  function handleSort(key) {
    if (sortKey === key) {
      sortAsc = !sortAsc;
    } else {
      sortKey = key;
      sortAsc = true;
    }
  }

  function extractLatestPcs(rawData) {
    const pcMap = new Map();
    rawData.forEach((pc) => {
      const serial = pc["Serial Number"];
      if (serial) pcMap.set(serial, pc);
    });
    return Array.from(pcMap.values());
  }

  onMount(async () => {
    try {
      const response = await fetch(API_URL);
      pcs = await response.json();
      loading = false;
    } catch (error) {
      console.error("Failed to load PC data:", error);
      loading = false;
    }
  });
</script>

<main class="dashboard-container">
  <Header />

  {#if loading}
    <div class="loading-state">
      <div class="spinner"></div>
      <p>Syncing with cloud database...</p>
    </div>
  {:else}
    <Controls
      bind:searchQuery
      bind:selectedLocation
      {locations}
      filteredCount={filteredPcs.length}
      totalCount={cleanPcs.length}
    />

    <div class="card table-card">
      <DeviceTable
        pcs={paginatedPcs}
        {sortKey}
        {sortAsc}
        onSort={handleSort}
        onSelect={(pc) => (selectedPc = pc)}
      />

      <Pagination
        {currentPage}
        {totalPages}
        onPrev={() => currentPage--}
        onNext={() => currentPage++}
      />
    </div>

    <DeviceDetail {selectedPc} onClose={() => (selectedPc = null)} />
  {/if}
</main>